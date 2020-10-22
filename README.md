# Maps - Admin panel for bybrisk Maps
This is a golang application. It uses docker for containerization and deployed on AWS EC2. Go.Mod and Go.Sum are used for dependency mannagement and docker is used for deployement.

## Getting Started with Maps
Map uses mongoDB for the databse and AWS S3 for the storage. Redis is used for cashing memory. shashankMongo repo handles the 
database connection and middlework repo handles the calculation part of the program.

* Razerpay is used as payment gateway.
* PubNub is used for streaming of live geolocation of delivery agents. 
* Mapbox API is used for any map related calculation and rendering.

## How to run Maps locally on your machine

  * Clone the repository to your device (by default it will be saved as github.com/shashank404error/maps).
  * Install golang and other dependencies using go get or by simply downloading go.mod file
  * Prefrably you need to create a mongodb account and get your credential from there.

 ## You can also get started by docker run shashank404error1/bybrisk-map 
                            
## Dependencies
   * **shashank404error/middlework**  - Handles the database connection and queries.
   * **shashank404error/shashankMongo** - Handles the calculation part of the application.

   
    You can view all the other dependencies in go.mod file
    

## Built With
Visual Studio Code

## Author
Shashank P. Sharma

