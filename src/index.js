/* @flow */
import React from 'react';

export default class Component extends React.Component {
  sum(a :number, b: number) :number {
    return a + b;
  }

  render() :any {
    return (
      <div className="wrapper">
         The sum of 13 and 37 would be {this.sum(13, 37)}
      </div>
    );
  }
}