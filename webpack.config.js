/* eslint-env node */
var path = require('path');
var webpack = require('webpack');
var NyanProgressPlugin = require('nyan-progress-webpack-plugin');

const env = process.env.NODE_ENV || 'development';

module.exports = {
  devtool: 'eval',
  entry: [
    './src/index'
  ],

  output: {
    path: path.join(__dirname, './build/'),
    filename: 'bundle.js',
    publicPath: '/'
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new NyanProgressPlugin(),
  ],

  module: {
    loaders: [{
      test: /\.js$/,
      loaders: ['babel'],
      include: path.join(__dirname, 'src')
    },
    {
      test: /\.css$/,
      loader: 'style-loader!css-loader'
    }]
  }
};