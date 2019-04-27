## riskmanagement (prototype)

*riskmanagement* is a small app (prototype) that helps you create risk matrix to manage the risks in your projects. 
For each project you have this app will create a risk matrix image where you can specify the risks of your 
project. Because this app is just a prototype is just storing all the data in memory that means if you stop the 
GO built-in server you will lose all the data; however, if you decide to use the app seriously I suggest you to use a 
different storage, like a sequel database. Also, every time you stop the GO built-in server you will not only lose the 
data, but the risk matrix images will be removed from disk as well—so keep in mind this. The design pattern I used to 
create this app is DDD. So it is easy to change things in the app (pug-and-play style) so you can improve or expand the
scope of the app. 

**NOTE**
This app is just a prototype of a more robust solution; however, I thought this might be interesting for someone so 
I decided to share part of the code.   

![example][/sample_img.png]

## Install
If you have go installed in your PC simply run from the app root directory:
```
$ go run cmd/server/main.go 
```

## Common Application Directories

### `/media`

*/media* contains all the risk matrix images created with the app.

## Built With

* go version go1.11.5 linux/amd64

## Contributing

## Authors 
* Otto Schuldt - *Initial work*

## TODO

* tests
* use a real storage like a sequel database to store the risk matrix data
* catch more errors 

## License

This project is licensed under the MIT License.

[sample_img.png]: sample_img.jpg "Title"