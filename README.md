React Component Boilerplate
===========================
[![Build Status](https://travis-ci.org/olegakbarov/react-component-boilerplate.svg?branch=master)](https://travis-ci.org/olegakbarov/react-component-boilerplate)
[![Dependency Status](https://img.shields.io/david/olegakbarov/react-component-boilerplate.svg)](https://david-dm.org/olegakbarov/react-component-boilerplate)
[![devDependency Status](https://img.shields.io/david/dev/strongloop/express.svg?maxAge=2592000)](https://david-dm.org/olegakbarov/react-component-boilerplate?dev=true)

Minimal [React](https://facebook.github.io/react/) component boilerplate with [Babel 6](http://babeljs.io/), [Webpack](https://webpack.github.io/), hot module replacement via [babel-plugin-react-transform](https://github.com/gaearon/babel-plugin-react-transform), [Flow](http://flowtype.org/), tests with [Tape](https://github.com/substack/tape) and `eslint`-friendly.

Inspired by Dan Abramov's [library boilerplate](https://github.com/gaearon/library-boilerplate)

### How-to

`$ npm install`

`$ npm run dev`

navigate to `localhost:8080`

### Considerations

You might want to tune `.babelrc` and `webpack.config` based on your needs.


### Why's

- Why not to use npm scripts over webpack CLI to run a dev server? — [that's why](https://github.com/webpack/webpack-dev-server/issues/106)

- Why so fancy import works `import Component from Component`? — checkout the aliases in webpack config.

### Anything?

Feedback appreciated!
