import React, {useEffect, useState} from 'react';
import './App.css';
import Login from "./components/Login/Login";
import {
  BrowserRouter as Router,
  Route,
  Switch,
} from "react-router-dom";
import Dashboard from "./components/Dashboard/Dashboard";
import {createMuiTheme, ThemeProvider} from "@material-ui/core/styles";
import {deepOrange, deepPurple, lightBlue, orange} from "@material-ui/core/colors";
import CssBaseline from "@material-ui/core/CssBaseline";
import AuthService from "./services/auth/auth";
import PolicyService from "./services/policy/policy";
import CircularProgress from "@material-ui/core/CircularProgress";


function App() {
  let isDarkTheme = false
  if (localStorage.getItem("theme") === "dark"){
    isDarkTheme = true
  }
  const [darkState, setDarkState] = useState(isDarkTheme);
  const [authExists, setAuthExists] = useState(AuthService.isAValidToken());
  const [currentPolicy, setPolicy] = useState(undefined);
  const [LoginMessage, setLoginMessage] = useState("");

  const updatePolicy = () =>{
    PolicyService.Get().then(response => {
      setPolicy(response)
    })

  }
  useEffect(() => {
    if (!authExists){
    AuthService.logAnonymousUser().then(() => {
          setAuthExists(true)
         updatePolicy()
        })
    } else if (!currentPolicy){
       updatePolicy()
    }
    const policyUpdateInterval = setInterval(() => {
      updatePolicy()
    }, 60000);
    return () => clearInterval(policyUpdateInterval);
  }, []);


  const palletType = darkState ? "dark" : "light";
  const mainPrimaryColor = darkState ? orange[500] : lightBlue[500];
  const mainSecondaryColor = darkState ? deepOrange[900] : deepPurple[500];
  const darkTheme = createMuiTheme({
    palette: {
      type: palletType,
      primary: {
        main: mainPrimaryColor
      },
      secondary: {
        main: mainSecondaryColor
      }
    }
  });

  return (
    <div className="App">
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />
        {authExists && currentPolicy ?
            <Router>
              <Switch>
                <Route exact path="/login" render={(props) => (
                    <Login {...props} message={LoginMessage} setMessage={setLoginMessage} setDarkState={setDarkState} darkState={darkState}/>
                )} />
                <Route path="/"   render={(props) => (
                    <Dashboard {...props} currentPolicy={currentPolicy} setLoginMessage={setLoginMessage} setDarkState={setDarkState} darkState={darkState}/>
                )} />
              </Switch>
            </Router>
            :
            <CircularProgress />
        }

    </ThemeProvider>
    </div>

  );
}

export default App;
