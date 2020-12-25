INSTALLATION WITH DOCKER
------------

### Prepare settings for DB connection 

Create file .env in project base directory. Copy the .env.example file to a local .env

Add your application configuration to a .env or use default congiguration for docker-containers.

### Build application with docker-compose 

Run docker-compose command:

~~~
docker-compose up --build
~~~

Now you should be able to access the application through the following URL:

~~~
http://localhost:85/
~~~

Adminer for management database through the following URL:

~~~
http://localhost:90/
~~~

