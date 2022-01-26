import { common, green, grey } from '@material-ui/core/colors'
import { createMuiTheme } from '@material-ui/core/styles'
import { darken } from '@material-ui/core/styles/colorManipulator'
import { ThemeOptions } from '@material-ui/core/styles/createMuiTheme'
import spacing from '@material-ui/core/styles/spacing'

declare module '@material-ui/core/styles/createPalette' {
  interface ListStatus {
    background: string
    color: string
  }

  interface PaletteOptions {
    success: PaletteColorOptions
    warning: PaletteColorOptions
    listPendingStatus: ListStatus
    listCompletedStatus: ListStatus
  }

  interface TypeBackground {
    appBar: string
  }
}

declare module '@material-ui/core/styles/createTypography' {
  type AdditionalThemeStyle = 'body1Next' | 'body2Next'
  interface TypographyStyleOptions {
    fontWeightLight?: number
    fontWeightMedium?: number
    fontWeightRegular?: number
    marginLeft?: string
  }

  interface TypographyOptions
    extends Partial<
      Record<AdditionalThemeStyle, TypographyStyleOptions> & FontStyleOptions
    > {}
}

const mainTheme: ThemeOptions = {
  props: {
    MuiGrid: {
      spacing: (spacing.unit * 3) as any as Required<
        Required<Required<ThemeOptions>['props']>['MuiGrid']
      >['spacing'],
    },
    MuiCardHeader: {
      titleTypographyProps: { color: 'secondary' },
    },
  },
  palette: {
    action: {
      hoverOpacity: 0.3,
    },
    primary: {
      light: '#E5F1FF',
      main: '#3c40c6',
      contrastText: '#fff',
    },
    secondary: {
      main: '#3d5170',
    },
    success: {
      light: '#e8faf1',
      main: green.A700,
      dark: green['700'],
      contrastText: common.white,
    },
    warning: {
      light: '#FFFBF1',
      main: '#fff6b6',
      contrastText: '#fad27a',
    },
    error: {
      light: '#ffdada',
      main: '#f44336',
      dark: '#d32f2f',
      contrastText: '#fff',
    },
    background: {
      default: '#f5f6f8',
      appBar: '#3c40c6',
    },
    text: {
      primary: darken(grey['A700'], 0.7),
      secondary: '#818ea3',
    },
    listPendingStatus: {
      background: '#fef7e5',
      color: '#fecb4c',
    },
    listCompletedStatus: {
      background: '#e9faf2',
      color: '#4ed495',
    },
  },
  shape: {
    borderRadius: spacing.unit,
  },
  overrides: {
    MuiButton: {
      root: {
        borderRadius: spacing.unit / 2,
        textTransform: 'none',
      },
      sizeLarge: {
        padding: undefined,
        fontSize: undefined,
        paddingTop: spacing.unit,
        paddingBottom: spacing.unit,
        paddingLeft: spacing.unit * 5,
        paddingRight: spacing.unit * 5,
      },
    },
    MuiTableCell: {
      body: {
        fontSize: '1rem',
      },
      head: {
        fontSize: '1rem',
        fontWeight: 400,
      },
    },
    MuiCardHeader: {
      root: {
        borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
      },
      action: {
        marginTop: -2,
        marginRight: 0,
        '& >*': {
          marginLeft: spacing.unit * 2,
        },
      },
      subheader: {
        marginTop: spacing.unit * 0.5,
      },
    },
  },
  typography: {
    useNextVariants: true,
    fontFamily: [
      '-apple-system',
      'BlinkMacSystemFont',
      'Roboto',
      'Helvetica',
      'Arial',
      'sans-serif',
    ].join(','),
    button: {
      textTransform: 'none',
      fontSize: '1.2em',
    },
    body1: {
      fontSize: '1.0rem',
      fontWeight: 400,
      lineHeight: '1.46429em',
      color: 'rgba(0, 0, 0, 0.87)',
      letterSpacing: -0.4,
    },
    body2: {
      fontSize: '1.0rem',
      fontWeight: 500,
      lineHeight: '1.71429em',
      color: 'rgba(0, 0, 0, 0.87)',
      letterSpacing: -0.4,
    },
    body1Next: {
      color: 'rgb(29, 29, 29)',
      fontWeight: 400,
      fontSize: '1rem',
      lineHeight: 1.5,
      letterSpacing: -0.4,
    },
    body2Next: {
      color: 'rgb(29, 29, 29)',
      fontWeight: 400,
      fontSize: '0.875rem',
      lineHeight: 1.5,
      letterSpacing: -0.4,
    },
    display1: {
      color: '#818ea3',
      fontSize: '2.125rem',
      fontWeight: 400,
      lineHeight: '1.20588em',
      letterSpacing: -0.4,
    },
    display2: {
      color: '#818ea3',
      fontSize: '2.8125rem',
      fontWeight: 400,
      lineHeight: '1.13333em',
      marginLeft: '-.02em',
      letterSpacing: -0.4,
    },
    display3: {
      color: '#818ea3',
      fontSize: '3.5rem',
      fontWeight: 400,
      lineHeight: '1.30357em',
      marginLeft: '-.02em',
      letterSpacing: -0.4,
    },
    display4: {
      fontSize: 14,
      fontWeightLight: 300,
      fontWeightMedium: 500,
      fontWeightRegular: 400,
      letterSpacing: -0.4,
    },
    h1: {
      color: 'rgb(29, 29, 29)',
      fontSize: '6rem',
      fontWeight: 300,
      lineHeight: 1,
    },
    h2: {
      color: 'rgb(29, 29, 29)',
      fontSize: '3.75rem',
      fontWeight: 300,
      lineHeight: 1,
    },
    h3: {
      color: 'rgb(29, 29, 29)',
      fontSize: '3rem',
      fontWeight: 400,
      lineHeight: 1.04,
    },
    h4: {
      color: 'rgb(29, 29, 29)',
      fontSize: '2.125rem',
      fontWeight: 400,
      lineHeight: 1.17,
    },
    h5: {
      color: 'rgb(29, 29, 29)',
      fontSize: '1.5rem',
      fontWeight: 400,
      lineHeight: 1.33,
      letterSpacing: -0.4,
    },
    h6: {
      fontSize: '0.8rem',
      fontWeight: 450,
      lineHeight: '1.71429em',
      color: 'rgba(0, 0, 0, 0.87)',
      letterSpacing: -0.4,
    },
    subheading: {
      color: 'rgb(29, 29, 29)',
      fontSize: '1rem',
      fontWeight: 400,
      lineHeight: '1.5em',
      letterSpacing: -0.4,
    },
    subtitle1: {
      color: 'rgb(29, 29, 29)',
      fontSize: '1rem',
      fontWeight: 400,
      lineHeight: 1.75,
      letterSpacing: -0.4,
    },
    subtitle2: {
      color: 'rgb(29, 29, 29)',
      fontSize: '0.875rem',
      fontWeight: 500,
      lineHeight: 1.57,
      letterSpacing: -0.4,
    },
  },
  shadows: [
    'none',
    '0px 1px 3px 0px rgba(0, 0, 0, 0.1),0px 1px 1px 0px rgba(0, 0, 0, 0.04),0px 2px 1px -1px rgba(0, 0, 0, 0.02)',
    '0px 1px 5px 0px rgba(0, 0, 0, 0.1),0px 2px 2px 0px rgba(0, 0, 0, 0.04),0px 3px 1px -2px rgba(0, 0, 0, 0.02)',
    '0px 1px 8px 0px rgba(0, 0, 0, 0.1),0px 3px 4px 0px rgba(0, 0, 0, 0.04),0px 3px 3px -2px rgba(0, 0, 0, 0.02)',
    '0px 2px 4px -1px rgba(0, 0, 0, 0.1),0px 4px 5px 0px rgba(0, 0, 0, 0.04),0px 1px 10px 0px rgba(0, 0, 0, 0.02)',
    '0px 3px 5px -1px rgba(0, 0, 0, 0.1),0px 5px 8px 0px rgba(0, 0, 0, 0.04),0px 1px 14px 0px rgba(0, 0, 0, 0.02)',
    '0px 3px 5px -1px rgba(0, 0, 0, 0.1),0px 6px 10px 0px rgba(0, 0, 0, 0.04),0px 1px 18px 0px rgba(0, 0, 0, 0.02)',
    '0px 4px 5px -2px rgba(0, 0, 0, 0.1),0px 7px 10px 1px rgba(0, 0, 0, 0.04),0px 2px 16px 1px rgba(0, 0, 0, 0.02)',
    '0px 5px 5px -3px rgba(0, 0, 0, 0.1),0px 8px 10px 1px rgba(0, 0, 0, 0.04),0px 3px 14px 2px rgba(0, 0, 0, 0.02)',
    '0px 5px 6px -3px rgba(0, 0, 0, 0.1),0px 9px 12px 1px rgba(0, 0, 0, 0.04),0px 3px 16px 2px rgba(0, 0, 0, 0.02)',
    '0px 6px 6px -3px rgba(0, 0, 0, 0.1),0px 10px 14px 1px rgba(0, 0, 0, 0.04),0px 4px 18px 3px rgba(0, 0, 0, 0.02)',
    '0px 6px 7px -4px rgba(0, 0, 0, 0.1),0px 11px 15px 1px rgba(0, 0, 0, 0.04),0px 4px 20px 3px rgba(0, 0, 0, 0.02)',
    '0px 7px 8px -4px rgba(0, 0, 0, 0.1),0px 12px 17px 2px rgba(0, 0, 0, 0.04),0px 5px 22px 4px rgba(0, 0, 0, 0.02)',
    '0px 7px 8px -4px rgba(0, 0, 0, 0.1),0px 13px 19px 2px rgba(0, 0, 0, 0.04),0px 5px 24px 4px rgba(0, 0, 0, 0.02)',
    '0px 7px 9px -4px rgba(0, 0, 0, 0.1),0px 14px 21px 2px rgba(0, 0, 0, 0.04),0px 5px 26px 4px rgba(0, 0, 0, 0.02)',
    '0px 8px 9px -5px rgba(0, 0, 0, 0.1),0px 15px 22px 2px rgba(0, 0, 0, 0.04),0px 6px 28px 5px rgba(0, 0, 0, 0.02)',
    '0px 8px 10px -5px rgba(0, 0, 0, 0.1),0px 16px 24px 2px rgba(0, 0, 0, 0.04),0px 6px 30px 5px rgba(0, 0, 0, 0.02)',
    '0px 8px 11px -5px rgba(0, 0, 0, 0.1),0px 17px 26px 2px rgba(0, 0, 0, 0.04),0px 6px 32px 5px rgba(0, 0, 0, 0.02)',
    '0px 9px 11px -5px rgba(0, 0, 0, 0.1),0px 18px 28px 2px rgba(0, 0, 0, 0.04),0px 7px 34px 6px rgba(0, 0, 0, 0.02)',
    '0px 9px 12px -6px rgba(0, 0, 0, 0.1),0px 19px 29px 2px rgba(0, 0, 0, 0.04),0px 7px 36px 6px rgba(0, 0, 0, 0.02)',
    '0px 10px 13px -6px rgba(0, 0, 0, 0.1),0px 20px 31px 3px rgba(0, 0, 0, 0.04),0px 8px 38px 7px rgba(0, 0, 0, 0.02)',
    '0px 10px 13px -6px rgba(0, 0, 0, 0.1),0px 21px 33px 3px rgba(0, 0, 0, 0.04),0px 8px 40px 7px rgba(0, 0, 0, 0.02)',
    '0px 10px 14px -6px rgba(0, 0, 0, 0.1),0px 22px 35px 3px rgba(0, 0, 0, 0.04),0px 8px 42px 7px rgba(0, 0, 0, 0.02)',
    '0px 11px 14px -7px rgba(0, 0, 0, 0.1),0px 23px 36px 3px rgba(0, 0, 0, 0.04),0px 9px 44px 8px rgba(0, 0, 0, 0.02)',
    '0px 11px 15px -7px rgba(0, 0, 0, 0.1),0px 24px 38px 3px rgba(0, 0, 0, 0.04),0px 9px 46px 8px rgba(0, 0, 0, 0.02)',
  ],
}

export const theme = createMuiTheme(mainTheme)
