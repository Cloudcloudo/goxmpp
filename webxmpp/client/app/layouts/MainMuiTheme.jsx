import getMuiTheme from 'material-ui/styles/getMuiTheme';

import { basicColors } from './Colors';

export const muiTheme = getMuiTheme({
  palette:{
    primary2Color: basicColors.cornFlower
  },
  textField: {
    textColor: basicColors.mineShaft,
    floatingLabelColor: basicColors.mineShaft,
    focusColor: basicColors.mariner,
    errorColor: basicColors.crimson,
    borderColor: basicColors.mineShaft
  },
  raisedButton: {
    primaryTextColor: '#ffffff',
    primaryColor: basicColors.cornFlower,
  },
  checkbox: {
    boxColor: basicColors.mineShaft,
    checkedColor: basicColors.mariner,
    requiredColor: basicColors.crimson,
    labelColor: basicColors.mineShaft
  },
  menuItem: {
    selectedTextColor: basicColors.mariner
  },
  datePicker: {
    color: basicColors.biscay,
    textColor: 'white',
    selectColor: basicColors.marine,
    selectTextColor: 'white',
    headerColor: basicColors.biscay,
  },
  // flatButton is in the datePicker
  flatButton: {
    primaryTextColor: basicColors.biscay,
  }
});
