import React from 'react';
import { configure, storiesOf, action, linkTo } from '@kadira/storybook';

import StaffjoyButton from '../components/StaffjoyButton';

storiesOf('StaffjoyButton', StaffjoyButton)
  .add('default', () => {
    return <StaffjoyButton>Hello world</StaffjoyButton>;
  });
