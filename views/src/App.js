import React from 'react';
import './App.css';
import Login from "./components/Login";
import {
  BrowserRouter as Router,
  Route,
  Switch,
  Link, withRouter
} from "react-router-dom";



function App() {

  return (
    <div className="App">
      <Router>
        <Switch>
          <Route exact path="/login" component={Login} />
        <Route exact path="/"> </Route>
        </Switch>
      </Router>
    </div>
  );
}

export default App;
