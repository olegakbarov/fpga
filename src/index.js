/* @flow */
import React from 'react';

export default class Component extends React.Component {
  sum(a :number, b: number) :number {
    return a + b;
  }

  render() :any {
    return (
      <div>
         The sum of 1337 and 7331 would be {this.sum(1337, 7331)}
      </div>
    );
  }
}