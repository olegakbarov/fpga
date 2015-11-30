/* eslint-env node */
var path = require('path');
var webpack = require('webpack');
var NyanProgressPlugin = require('nyan-progress-webpack-plugin');

const env = process.env.NODE_ENV || 'development';

module.exports = {
  devtool: 'eval',
  entry: './example/index.js',

  output: {
    path: path.join(__dirname, '.'),
    filename: 'bundle.js',
    publicPath: '/'
  },

  resolve: {
    alias: {
      'Component': path.join(__dirname, '../src')
    },
    extensions: ['', '.js']
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new NyanProgressPlugin(),
  ],

  module: {
    loaders: [{
      test: /\.js$/,
      loaders: ['babel'],
      include: [ path.join(__dirname, '../src'), path.join(__dirname, '.') ]
    },
    {
      test: /\.css$/,
      loader: 'style-loader!css-loader'
    }]
  }
};