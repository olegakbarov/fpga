/* eslint-env node */
var webpack = require('webpack');
var WebpackDevServer = require('webpack-dev-server');
var config = require('./webpack.config.js');

new WebpackDevServer(webpack(config), {
  hot: true,
  publicPath: config.output.publicPath,
  historyApiFallback: true,
  stats: {
    colors: true
  }
}).listen(8080, 'localhost', () => console.log('done'));