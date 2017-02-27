# Staffjoy Web Assets

## What is it?

Hosted Javascript and CSS files, for use across all Staffjoy websites.

## Using the assets

In the header of your html page, include a reference to `/bundle.js` and `styles.css` which are hosted on root of this domain.

### What about fonts?

For now, fonts will not be included. To make sure you are using the correct version of Open Sans, include this line in your header.

```
<link href='https://fonts.googleapis.com/css?family=Open+Sans:400,300,300italic,400italic,600,700,800' rel='stylesheet' type='text/css'>
```

## Developers

Below are the scripts needed for running and modifying the products of this microservice.

### Getting Started

```npm install```

### dev
```npm start```

http://localhost:8080 will have ```bundle.js``` and ```styles.css``` available for consumption. If you modify any of the shared resources, they will automatically be rebuilt.

### production
```npm run build```

The folder ```web_assets/dist/``` will be populated with the ```bundle.js``` and ```styles.css``` files.

