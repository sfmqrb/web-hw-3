import React, { Component } from "react";
import jwtDecode from "jwt-decode";
import { Route, Redirect, Switch } from "react-router-dom";
import "react-toastify/dist/ReactToastify.css";
import { ToastContainer } from "react-toastify";
import Notes from "./components/notes";
import NoteForm from "./components/noteForm";
import NotFound from "./components/notFound";
import NavBar from "./components/navBar";
import LoginForm from "./components/loginForm";
import RegisterForm from "./components/registerForm";
import Logout from "./components/logout";
import "./App.css";

class App extends Component {
  state = {};

  // backend related
  // componentDidMount() {
  // try{
  //   const jwt = localStorage.getItem("jwt");
  //   const user = jwtDecode(jwt)
  //   this.setState({ user });
  // }
  // catch(e){
  // navigate to login page
  // }
  // }

  render() {
    return (
      <React.Fragment>
        <ToastContainer />
        <NavBar />
        <main className="container">
          <Switch>
            <Route path="/register" component={RegisterForm} />
            <Route path="/login" component={LoginForm} />
            <Route path="/logout" component={Logout} />
            <Route path="/notes/:id" component={NoteForm} />
            <Route path="/notes" component={Notes} />
            <Route path="/not-found" component={NotFound} />
            <Redirect from="/" exact to="/notes" />
            <Redirect to="/not-found" />
          </Switch>
        </main>
      </React.Fragment>
    );
  }
}

export default App;
