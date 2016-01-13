import React from 'react';
import test from 'tape';
import { shallow } from 'enzyme';
import Component from '../src/index';

test('correct class', assert => {
  const msg = 'should render component with correct class';

  let expected = `<div class="wrapper">The sum of 13 and 37 would be 50</div>`;

  const $ = shallow(<Component />);
  const output = $.html();

  assert.equal(output, expected, msg);

  assert.end();
});