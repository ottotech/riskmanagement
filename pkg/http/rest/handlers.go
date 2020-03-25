package rest

import (
	"encoding/json"
	"fmt"
	"github.com/ottotech/riskmanagement/pkg/adding"
	"github.com/ottotech/riskmanagement/pkg/config"
	"github.com/ottotech/riskmanagement/pkg/deleting"
	"github.com/ottotech/riskmanagement/pkg/draw"
	"github.com/ottotech/riskmanagement/pkg/listing"
	"github.com/ottotech/riskmanagement/pkg/updating"
	"github.com/ottotech/riskmanagement/pkg/utils"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type App struct {
	ListMatrix       *ListMatrix
	AddMatrix        *AddMatrix
	GetMatrix        *GetMatrix
	Media            *Media
	AddRisk          *AddRisk
	DeleteRisk       *DeleteRisk
	DeleteRiskMatrix *DeleteRiskMatrix
	AddMediaPath     *AddMediaPath
}

type ListMatrix struct {
}

func (h *ListMatrix) Handler(s listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		list := s.GetAllRiskMatrix()
		utils.RenderTemplate(w, "templates/list.gohtml", list)
		return
	}
}

type AddMatrix struct {
}

func (h *AddMatrix) Handler(a adding.Service, l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			utils.RenderTemplate(w, "templates/add.gohtml", nil)
			return
		}
		var err error
		name := r.PostFormValue("project")
		if name == "" {
			utils.RenderTemplate(w, "templates/add.gohtml", "Error: You need to specify the project name.")
			return
		}
		t := time.Now().Format("02_01_2006_03_04_05")
		filename := fmt.Sprintf("%v.png", t)
		rm := adding.RiskMatrix{Project: name, Path: filename}
		err = a.AddRiskMatrix(rm)
		if err != nil {
			config.Logger.Println(err)
			utils.RenderTemplate(w, "templates/add.gohtml", fmt.Sprintf("There was an internal error."))
			return
		}
		newRm, err := l.GetRiskMatrixByPath(filename)
		if err != nil {
			config.Logger.Println(err)
			utils.RenderTemplate(w, "templates/add.gohtml", fmt.Sprintf("There was an internal error."))
			return
		}
		// TODO: What happens if we add the risk matrix data but we cannot draw the risk matrix
		// TODO: for some reason?
		mediaPath, err := l.GetMediaPath()
		if err != nil {
			config.Logger.Println(err)
			utils.RenderTemplate(w, "templates/add.gohtml", err.Error())
			return
		}
		pathToDraw := filepath.Join(mediaPath, filename)
		err = draw.RiskMatrixDrawer(pathToDraw, newRm, []adding.Risk{})
		if err != nil {
			config.Logger.Println(err)
			utils.RenderTemplate(w, "templates/add.gohtml", err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

type DeleteRiskMatrix struct {
}

func (h *DeleteRiskMatrix) Handler(d deleting.Service, l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get riskID from request
		id, err := strconv.Atoi(r.PostFormValue("risk_matrix_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// before deleting the matrix, let's instantiate it
		riskMatrix, err := l.GetRiskMatrix(id)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete the risk matrix
		err = d.DeleteRiskMatrix(id)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete image
		mediaPath, err := l.GetMediaPath()
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = os.Remove(filepath.Join(mediaPath, riskMatrix.Path))
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: What happens when we delete the risk matrix but for some reason we are not able to remove
		// TODO: all the risks?
		// delete all risks from matrix
		risks := l.GetAllRisks(riskMatrix.ID)
		risksIDs := make([]string, 0)
		for _, r := range risks {
			risksIDs = append(risksIDs, r.ID)
		}
		err = d.DeleteRisk(risksIDs...)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// if all goes well we return response 200
		w.WriteHeader(http.StatusOK)
	}
}

type AddRisk struct {
}

func (h *AddRisk) Handler(a adding.Service, l listing.Service, u updating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get all risks from request
		var risks []adding.Risk
		postData := r.PostFormValue("data")
		err := json.Unmarshal([]byte(postData), &risks)
		if err != nil {
			config.Logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// anonymous func that returns the risk classification
		riskClassifier := func(r adding.Risk) (string, bool) {
			if r.Probability == 3 && r.Impact == 1 {
				return "medium", true
			}
			if r.Probability == 3 && r.Impact == 2 {
				return "high", true
			}
			if r.Probability == 3 && r.Impact == 3 {
				return "high", true
			}
			if r.Probability == 2 && r.Impact == 1 {
				return "low", true
			}
			if r.Probability == 2 && r.Impact == 2 {
				return "medium", true
			}
			if r.Probability == 2 && r.Impact == 3 {
				return "high", true
			}
			if r.Probability == 1 && r.Impact == 1 {
				return "low", true
			}
			if r.Probability == 1 && r.Impact == 2 {
				return "low", true
			}
			if r.Probability == 1 && r.Impact == 3 {
				return "medium", true
			}
			return "", false
		}

		// we need to set the risk classification for each risk
		for i, r := range risks {
			c, ok := riskClassifier(r)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				msg := fmt.Sprintf("Probability and impact numbers should be numbers between 1 and 3 on risk (%v)", r.Name)
				_, _ = w.Write([]byte(msg))
				return
			}
			risks[i].Classification = c
		}

		// classify new risks
		var newRisks []adding.Risk
		riskMatrixID := risks[0].RiskMatrixID
		preexistingRisks := l.GetAllRisks(riskMatrixID)
		for _, r := range risks {
			exists := false
			for _, pr := range preexistingRisks {
				if r.Name == pr.Name {
					exists = true
				}
			}
			if !exists {
				newRisks = append(newRisks, r)
			}
		}

		// if there are no new risks we send a response with status code 200
		if len(newRisks) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}
		// adding risks
		err = a.AddRisk(newRisks...)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get risk matrix
		riskMatrix, err := l.GetRiskMatrix(riskMatrixID)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// this anonymous func will count how many risks are per block in the risk matrix.
		// then it will get the number of risks from the block that has the most number of risks.
		// finally based on this max number this func will resize the matrix in order to show the
		// risks accordingly.
		riskMatrixResize := func(rm listing.RiskMatrix, r []adding.Risk) {
			var rb1, rb2, rb3, rb4, rb5, rb6, rb7, rb8, rb9 int // rb = risk matrix block
			for _, r := range risks {
				if r.Probability == 3 && r.Impact == 1 {
					rb1 += 1
				}
				if r.Probability == 3 && r.Impact == 2 {
					rb2 += 1
				}
				if r.Probability == 3 && r.Impact == 3 {
					rb3 += 1
				}
				if r.Probability == 2 && r.Impact == 1 {
					rb4 += 1
				}
				if r.Probability == 2 && r.Impact == 2 {
					rb5 += 1
				}
				if r.Probability == 2 && r.Impact == 3 {
					rb6 += 1
				}
				if r.Probability == 1 && r.Impact == 1 {
					rb7 += 1
				}
				if r.Probability == 1 && r.Impact == 2 {
					rb8 += 1
				}
				if r.Probability == 1 && r.Impact == 3 {
					rb9 += 1
				}
			}
			var max int
			for _, v := range []int{rb1, rb2, rb3, rb4, rb5, rb6, rb7, rb8, rb9} {
				if v > max {
					max = v
				}
			}
			// block size:   200
			// border width: 3*2 (top and bottom)
			// word height:  13
			// line spacing: word height + 2
			blockWritableSize := 200 - 30 // 30 is a dummy value to add more space while adding risk labels on image
			actualSize := 0
			for i := 1; i < max; i++ {
				actualSize += 15
			}
			// if the size required is greater than the block writable size
			// we need to find the right size where the risks can fit.
			// note that the size of the blocks should be a reminder of a number
			// that is a multiple of 3 because the risk matrix has 3 rows and 3 columns
			if actualSize > blockWritableSize {
				actualSize = actualSize + 30 // 30 is a dummy value to add more space while adding risk labels on image
				newBlockWidth := 0
				multiple := riskMatrix.Multiple
				for {
					x := multiple + (3 - multiple%3)
					if x > actualSize {
						newBlockWidth = x
						break
					}
					multiple = x
				}
				newImgWidth := newBlockWidth * 3
				riskMatrix.MatImgWidth = newImgWidth
				riskMatrix.MatImgHeight = newImgWidth
				riskMatrix.Multiple = newImgWidth / riskMatrix.MatNrCols

				// finally we need to update the risk matrix sizes in the storage
				err := u.UpdateRiskMatrixSize(riskMatrixID, newImgWidth)
				if err != nil {
					w.WriteHeader(http.StatusForbidden)
					_, _ = w.Write([]byte(fmt.Sprintf("While updating the size of the risk "+
						"matrix we encounter this err: %v", err.Error())))
					return
				}
			}
		}

		// resize risk matrix image if necessary
		riskMatrixResize(riskMatrix, risks)

		// TODO: What happens if we add the risks but we cannot draw the risk matrix again?
		mediaPath, err := l.GetMediaPath()
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pathToDraw := filepath.Join(mediaPath, riskMatrix.Path)
		// draw risk matrix again in order to add the new risks
		err = draw.RiskMatrixDrawer(pathToDraw, riskMatrix, risks)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// if all goes well we send a status 200
		w.WriteHeader(http.StatusOK)
	}
}

type DeleteRisk struct {
}

func (h *DeleteRisk) Handler(d deleting.Service, l listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only allow POST method
		if r.Method == http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// response will be json
		w.Header().Set("Content-Type", "application/json")

		// get riskID from request
		id := r.PostFormValue("risk_id")

		// we get the risk instance before deleting it to use it later
		risk, err := l.GetRisk(id)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete the risk
		err = d.DeleteRisk(id)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get riskMatrix of the deleted risk
		riskMatrix, err := l.GetRiskMatrix(risk.RiskMatrixID)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get all risks of the matrix
		var risks []adding.Risk
		for _, risk := range l.GetAllRisks(riskMatrix.ID) {
			r := adding.Risk{
				RiskMatrixID:   riskMatrix.ID,
				Name:           risk.Name,
				Probability:    risk.Probability,
				Impact:         risk.Impact,
				Classification: risk.Classification,
				Strategy:       risk.Strategy,
			}
			risks = append(risks, r)
		}

		// draw risk matrix again in order to not show the deleted risks
		mediaPath, err := l.GetMediaPath()
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pathToDraw := filepath.Join(mediaPath, riskMatrix.Path)
		err = draw.RiskMatrixDrawer(pathToDraw, riskMatrix, risks)
		if err != nil {
			config.Logger.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// if all goes well we return response 200
		w.WriteHeader(http.StatusOK)
	}
}

type GetMatrix struct {
}

func (h *GetMatrix) Handler(lister listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// clean path to get RiskMatrix ID
			p := path.Clean("/" + r.URL.Path)
			i := strings.Index(p[1:], "/") + 1
			tail := p[i+1:]

			// check if id is valid
			id, err := strconv.Atoi(tail)
			if err != nil {
				log.Println(err)
				http.Error(w, fmt.Sprintf("This Risk Matrix ID is not valid: %v.", tail), http.StatusBadRequest)
				return
			}

			// get risk matrix
			riskMatrix, err := lister.GetRiskMatrix(id)
			if err != nil {
				config.Logger.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// get all risks
			risks := lister.GetAllRisks(riskMatrix.ID)

			// Let's sort the risks by name.
			sort.Sort(listing.ByName(risks))

			// build context
			ctx := struct {
				RiskMatrix listing.RiskMatrix
				Risks      []listing.Risk
			}{
				riskMatrix,
				risks,
			}

			// render tpl
			utils.RenderTemplate(w, "templates/detail.gohtml", ctx)
			return
		}
	}
}

type Media struct {
}

func (h *Media) Handler(lister listing.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := "image/png"
		w.Header().Add("Content-Type", ct)

		idx := strings.LastIndex(r.URL.Path, "/")
		filename := r.URL.Path[idx+1:]
		mediaPath, err := lister.GetMediaPath()
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, http.StatusText(404), http.StatusNotFound)
				return
			}
			config.Logger.Println(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		fullPath := filepath.Join(mediaPath, filename)
		f, err := os.Open(fullPath)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		fi, err := f.Stat()
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, r, f.Name(), fi.ModTime(), f)
	})
}

type AddMediaPath struct {
}

func (h *AddMediaPath) Handler(adder adding.Service, lister listing.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ctx := struct {
				Error, MediaPath string
			}{
				Error: "",
			}
			if mediaPath, err := lister.GetMediaPath(); err == nil {
				ctx.MediaPath = mediaPath
			}
			utils.RenderTemplate(w, "templates/mediapath.gohtml", ctx)
			return
		}
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			ctx := struct {
				Error, MediaPath string
			}{}
			mediaPath := r.PostForm.Get("mediapath")
			if mediaPath == "" {
				ctx.Error = "You need to specify a valid path."
				utils.RenderTemplate(w, "templates/mediapath.gohtml", ctx)
				return
			}
			mediaPath = filepath.Clean(mediaPath)
			if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
				ctx.Error = "The path you provided does not exist."
				ctx.MediaPath = r.PostForm.Get("mediapath")
				utils.RenderTemplate(w, "templates/mediapath.gohtml", ctx)
				return
			}
			f, err := os.Open(mediaPath)
			if err != nil {
				config.Logger.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			defer func() {
				err := f.Close()
				if err != nil {
					config.Logger.Println(err)
				}
			}()
			fi, err := f.Stat()
			if err != nil {
				config.Logger.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			if !fi.IsDir() {
				ctx.Error = "You need to specify a valid path to a folder, not to a file."
				ctx.MediaPath = r.PostForm.Get("mediapath")
				utils.RenderTemplate(w, "templates/mediapath.gohtml", ctx)
				return
			}
			err = adder.SaveMediaPath(mediaPath)
			if err != nil {
				config.Logger.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}
