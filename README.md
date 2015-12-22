# React Component Boilerplate

Basic [React](https://facebook.github.io/react/) component boilerplate with [Babel6](http://babeljs.io/), [Webpack](https://webpack.github.io/), hot module replacement via [transform](https://github.com/gaearon/babel-plugin-react-transform) [Flow](http://flowtype.org/) and `eslint`-ready.

### How-to

`$ npm install`

`$ npm run dev`

navigate to `localhost:8080`

### Got error?

```
ERROR in ./example/root.js
Module build failed: TypeError: Plugin is not a function
```

Try `npm install babel-plugin-react-transform@beta --save-dev` or advice me how to allow `npm i` to install betas of packages :)

### Considerations

You might want to tune `.babelrc` and `webpack.config` based on your needs.


### Why's

- Why not to use npm scripts over webpack CLI to run a dev server? — [that's why](https://github.com/webpack/webpack-dev-server/issues/106)

- Why so fancy import works `import Component from Component`? — checkout the aliases in webpack config.
