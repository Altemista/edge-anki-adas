# Agular Frontend for mwc-twin Component
This project is generated with [yo angular generator](https://github.com/yeoman/generator-angular)
version 0.16.0.

## Build & development
The dependencies are not version controlled, hence install dependencies with the following commands:
```
npm install
bower install
```

Then run `grunt` for building and `grunt serve` for preview. (TODO: The preview does not yet work as the app currently assumes a base URI of "html"). Note that the `./dist` folder is checked into git for faster deployments of the micro service.

## Testing
Running `grunt test` will run the unit tests with karma. Unfortunately, there are not tests yet ;-)

## Tooling
The frontend development requires the following tools:
```
npm install -g bower 
npm install -g yo generator-angular generator-karma
```

There were some issues with the global install of karma, hence:
```
npm install grunt grunt-cli jit-grunt 
npm install karma-phantomjs-launcher karma-jasmine jasmine-core phantomjs-prebuilt
npm install grunt-karma karma
```

## References
[Getting Started](https://www.sitepoint.com/kickstart-your-angularjs-development-with-yeoman-grunt-and-bower/)
