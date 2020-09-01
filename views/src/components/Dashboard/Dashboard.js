import React, { useState } from "react";
import clsx from "clsx";
import { makeStyles } from "@material-ui/core/styles";
import CssBaseline from "@material-ui/core/CssBaseline";
import Switch from "@material-ui/core/Switch";
import Drawer from "@material-ui/core/Drawer";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import List from "@material-ui/core/List";
import Typography from "@material-ui/core/Typography";
import Divider from "@material-ui/core/Divider";
import IconButton from "@material-ui/core/IconButton";
import Container from "@material-ui/core/Container";
import Grid from "@material-ui/core/Grid";
import MenuIcon from "@material-ui/icons/Menu";
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import { adminListItems } from "./listItems";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import ScoreBoard from "../ScoreBoard/ScoreBoard";
import AuthService from "../../services/auth/auth";
import {Link, Route} from "react-router-dom";
import AssignmentIndIcon from '@material-ui/icons/AssignmentInd';
import CloseIcon from '@material-ui/icons/Close';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import {Alert} from "@material-ui/lab";
import { FullScreen, useFullScreenHandle } from "react-full-screen";
import Settings from "../Admin/Settings";
import Setup from "../Admin/Setup";
import BarChartIcon from "@material-ui/icons/BarChart";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";
import DetailsIcon from "@material-ui/icons/Details";
import CircularProgress from "@material-ui/core/CircularProgress";
import Button from "@material-ui/core/Button";
import DialogActions from "@material-ui/core/DialogActions";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";

const drawerWidth = 260;

const useStyles = makeStyles(theme => ({
  root: {
    display: "flex"
  },
  toolbar: {
    position: "relative",
    paddingRight: 24
  },
  toolbarIcon: {
    display: "flex",
    alignItems: "center",
    justifyContent: "flex-end",
    padding: "0 8px",
    ...theme.mixins.toolbar
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(["width", "margin"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen
    })
  },
  appBarShift: {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(["width", "margin"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen
    })
  },
  menuButton: {
    marginRight: 36
  },
  menuButtonHidden: {
    display: "none"
  },
  title: {
    flexGrow: 1
  },
  drawerPaper: {
    position: "relative",
    whiteSpace: "nowrap",
    width: drawerWidth,
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen
    })
  },
  drawerPaperClose: {
    overflowX: "hidden",
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen
    }),
    width: theme.spacing(7),
    [theme.breakpoints.up("sm")]: {
      width: theme.spacing(9)
    }
  },
  appBarSpacer: theme.mixins.toolbar,
  content: {
    flexGrow: 1,
    height: "100vh",
    overflow: "auto"
  },
  container: {
    paddingTop: theme.spacing(4),
    paddingBottom: theme.spacing(4)
  },
  paper: {
    padding: theme.spacing(2),
    display: "flex",
    overflow: "auto",
    flexDirection: "column"
  },
  fixedHeight: {
    height: '85vh'
  },
  fullSizeHeight: {
    height: '100vh'
  }
}));

export default function Dashboard(props) {

  const [open, setOpen] = useState(false);
  const [Title, setTitle] = useState("ScoreBoard")
  const [alert, setAlert] = useState({message: "", severity: "", loader: false, openDialogue: false})
  const setDarkState = props.setDarkState
  const darkState = props.darkState
  const classes = useStyles();

  const handleOpenDialogue = () =>{
    setAlert(prevState => {return {...prevState, openDialogue: true}})
  }

  const handleCloseAlertDialogue = () =>{
    setAlert({message: "", severity: "", loader: false, openDialogue: false})
  }

  const handleCancelDialogue = () =>{
    setAlert(prevState => {return {...prevState, openDialogue: false}})
  }

  const handleSuccess = (msg = "Success") => {
    setAlert(prevState => {return { ...prevState, message: msg, severity: "success", loader: false}})
  }
  const handleLoading = () => {
    setAlert(prevState => {return { ...prevState, message: "", severity: "warning", loader: true}})
  }

  const cancelLoading = () => {
    setAlert(prevState => {return { ...prevState, message: "", severity: "", loader: false}})
  }


  const handleThemeChange = (e) => {
    if (e.target.checked){
      localStorage.setItem("theme", "dark")
    } else{
      localStorage.setItem("theme", "light")
    }
    setDarkState(e.target.checked);
  };

  const handleDrawerOpen = () => {
    setOpen(true);
  };
  const handleDrawerClose = () => {
    setOpen(false);
  };
  const fixedHeightPaper = clsx(classes.paper, classes.fixedHeight);
  const fullSizeHeightPaper = clsx(classes.paper, classes.fullSizeHeight);

  const errorSetter = (error) => {
    let resMessage =
        (error.response
            && error.response.data && ( error.response.data.message || error.response.data.error) ) ||
        error.toString();
    if (error.response.status === 403){
      props.history.push("/login");
      props.setLoginMessage("You need to Log-in in order to access this resource")
    } else {
      if (typeof resMessage != "string"){
        resMessage = "The request was invalid"
      }
      setAlert({message: resMessage, severity: "error", loader: false})
    }
  }


  const logout = () => {
    AuthService.logout()
    window.location.reload(true)
  }
  const handleFullScreen = useFullScreenHandle()
  return (

      <div className={classes.root}>
        <CssBaseline />
        <AppBar
          position="absolute"
          className={clsx(classes.appBar, open && classes.appBarShift)}
        >
          <Toolbar className={classes.toolbar}>
            <IconButton
              edge="start"
              color="inherit"
              aria-label="open drawer"
              onClick={handleDrawerOpen}
              className={clsx(
                classes.menuButton,
                open && classes.menuButtonHidden
              )}
            >
              <MenuIcon />
            </IconButton>
            <Typography
              component="h1"
              variant="h6"
              color="inherit"
              noWrap
              className={classes.title}
            >{Title}
            </Typography>
            <div style={{position: "absolute", "right": "10px"}}>
              {
                alert.message === "" && alert.loader &&
                <CircularProgress color="secondary" />
              }
              {alert.message && (
                  <Alert severity={alert.severity}
                         action={
                           <IconButton
                               aria-label="close"
                               color="inherit"
                               size="small"
                               onClick={() => {
                                 setAlert({message: "", severity: "", loader: false, openDialogue: false});
                               }}
                           >
                             <CloseIcon fontSize="inherit" />
                           </IconButton>
                         }
                  >{(() => {
                    const maxLen = 80
                    return (alert.message.length > maxLen) ?
                        <div>
                          {alert.message.substr(0, maxLen-1) + '...'}
                          <Button color="primary" onClick={handleOpenDialogue} style={{marginLeft:"20px"}} size="small">
                            Open Details
                          </Button>
                          <Dialog
                              open={alert.openDialogue}
                              onClose={handleCancelDialogue}
                              aria-labelledby="alert-dialog-title"
                              aria-describedby="alert-dialog-description"
                          >
                            <DialogTitle id="alert-dialog-title">Alert Details</DialogTitle>
                            <DialogContent>
                              <DialogContentText id="alert-dialog-description">
                                {alert.message.split("\n").map(text => {
                                  return (
                                      <Typography>
                                        {text}
                                      </Typography>
                                  )
                                })}
                              </DialogContentText>
                            </DialogContent>
                            <DialogActions>
                              <Button onClick={handleCancelDialogue} color="primary">
                                Close Dialogue
                              </Button>
                              <Button onClick={handleCloseAlertDialogue} color="primary" autoFocus>
                                Close Alert
                              </Button>
                            </DialogActions>
                          </Dialog>
                        </div>
                        : alert.message;
                  })()}</Alert>
              )}
            </div>
          </Toolbar>
        </AppBar>
        <Drawer
          variant="permanent"
          classes={{
            paper: clsx(classes.drawerPaper, !open && classes.drawerPaperClose)
          }}
          open={open}>
          <div className={classes.toolbarIcon}>
            <IconButton onClick={handleDrawerClose}>
              <ChevronLeftIcon />
            </IconButton>
          </div>
          <Divider/>
          <Switch checked={darkState} onChange={handleThemeChange} />
          <Divider />
          <List>
            <div>
              {!(AuthService.getCurrentRole() !== "black" && !props.currentPolicy["allow_to_see_points"]) &&
              <ListItem button component={Link} to="/ranks">
                <ListItemIcon>
                  <BarChartIcon/>
                </ListItemIcon>
                <ListItemText primary="Ranks" />
              </ListItem>
              }
              <ListItem button component={Link} to="/">
                <ListItemIcon>
                  <CheckCircleIcon />
                </ListItemIcon>
                <ListItemText primary="Status" />
              </ListItem>
              { (AuthService.getCurrentRole() === "blue" || AuthService.getCurrentRole() ==="black") &&
              <ListItem button component={Link} to="/details">
                <ListItemIcon>
                  <DetailsIcon />
                </ListItemIcon>
                <ListItemText primary="Details" />
              </ListItem>
              }
            </div>
          </List>
          {
            AuthService.getCurrentRole() === "black" &&
            <List>
              <Divider/>
                {adminListItems}
            </List>
          }
          <Divider/>
          {
            !AuthService.isAValidRole() ?
                <ListItem button component={Link} to="/login">
                  <ListItemIcon>
                    <AssignmentIndIcon />
                  </ListItemIcon>
                  <ListItemText primary="Sign In" />
                </ListItem>
                :
                <ListItem button onClick={logout}>
                  <ListItemIcon>
                    <ExitToAppIcon />
                  </ListItemIcon>
                  <ListItemText primary="Sign Out" />

                </ListItem>
          }
        </Drawer>
        <main className={classes.content}>
          <div className={classes.appBarSpacer} />
          <Container maxWidth="xl" className={classes.container}>
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <Route exact path={["/", "/ranks", "/details"]} render={() => (
                    <FullScreen handle={handleFullScreen}>
                        <ScoreBoard currentPolicy={props.currentPolicy} errorSetter={errorSetter} isDarkTheme={darkState} setTitle={setTitle} handleFullScreen={handleFullScreen} fullSizeHeightPaper={fullSizeHeightPaper} fixedHeightPaper={fixedHeightPaper}/>
                    </FullScreen>
                )} />
                <Route exact path="/settings" render={() => (
                    <Settings errorSetter={errorSetter} isDarkTheme={darkState} cancelLoading={cancelLoading} handleSuccess={handleSuccess} handleLoading={handleLoading} setAlert={setAlert} setTitle={setTitle} classesPaper={classes.paper}/>
                )} />
                <Route path="/setup" render={() => (
                    <Setup  errorSetter={errorSetter} isDarkTheme={darkState} cancelLoading={cancelLoading} handleSuccess={handleSuccess} handleLoading={handleLoading}  setAlert={setAlert} setTitle={setTitle} classesPaper={classes.paper}/>
                )} />

              </Grid>
            </Grid>
          </Container>
        </main>
      </div>

  );
}