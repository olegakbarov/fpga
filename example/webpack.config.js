/* eslint-env node */
/* eslint-global env */
const path = require('path');
const webpack = require('webpack');

module.exports = {
  devtool: 'cheap-inline-module-source-map',
  entry: [
    'webpack-dev-server/client?http://localhost:8080',
    'webpack/hot/only-dev-server',
    path.join(__dirname, './root')
  ],

  output: {
    path: path.join(__dirname, '.'),
    publicPath: '/assets/',
    filename: 'bundle.js'
  },

  resolve: {
    alias: {
      Component: path.join(__dirname, '../src')
    },
    extensions: ['', '.js']
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin()
  ],

  devServer: {
    contentBase: __dirname
  },

  module: {
    loaders: [{
      test: /\.js$/,
      loader: 'babel-loader',
      exclude: /node_modules/,
      include: [
        __dirname,
        path.resolve(__dirname, '../src')
      ]
    },
    {
      test: /\.css$/,
      loader: 'style-loader!css-loader'
    }]
  }
};
