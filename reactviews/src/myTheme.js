import blue from '@material-ui/core/colors/blue';
import indigo from '@material-ui/core/colors/indigo';
import cyan from '@material-ui/core/colors/cyan';
import grey from '@material-ui/core/colors/grey';
import white from '@material-ui/core/colors';
import black from '@material-ui/core/colors';

import { fade } from '@material-ui/core/styles/colorManipulator';
//import spacing from '@material-ui/core/styles';

const myTheme = {
    //spacing: spacing,
    fontFamily: 'Roboto, sans-serif',
    palette: {
        primary: blue,
        secondary: {main: '#2196F3'},
        primary1Color: blue[500],
        primary2Color: blue[700],
        primary3Color: grey[400],
        accent1Color: indigo[500],
        accent2Color: grey[100],
        accent3Color: grey[500],
        textColor: grey[900],
        alternateTextColor: white,
        canvasColor: white,
        borderColor: grey[300],
        disabledColor: fade(grey[800], 0.3),
        pickerHeaderColor: cyan[500],
        clockCircleColor:  fade(grey[800], 0.07),
        shadowColor: black,
    },  
};

export default myTheme;
