# www

Short story: Run this command **before committing**.

```
bash build.sh
```

## CSS

**DO NOT DIRECTLY MODIFY CSS**

CSS is being compiled from Sass, but we are bundling everything into 1 file, located in `assets/css/staffjoy.css`. To write new CSS, write it in the appropriate `.scss` file located in the `sass/` folder, then run ```gulp``` to build it. Gulp will watch for changes and rebuild, so you can leave it running in a separate window. Your SCSS changes will spur changes in `assets/css/staffjoy.css` - **be sure to commit these**.

## Building assets for Go

If you make changes to assets or templates, you must run the following scripts to package your assets for the binary:

```
go-bindata assets/...
```

the bindata package is old and ornery, so you may need to run `gofmt -s` (with the simplification flag) or linting on it.


