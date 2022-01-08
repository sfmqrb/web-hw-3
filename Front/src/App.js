import React, { Component } from "react";
import { Route, Redirect, Switch } from "react-router-dom";
import Notes from "./components/notes";
import NoteForm from "./components/noteForm";
import NotFound from "./components/notFound";
import NavBar from "./components/navBar";
import LoginForm from "./components/loginForm";
import RegisterForm from "./components/registerForm";
import "./App.css";

class App extends Component {
  render() {
    return (
      <React.Fragment>
        <NavBar />
        <main className="container">
          <Switch>
            <Route path="/register" component={RegisterForm} />
            <Route path="/login" component={LoginForm} />
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
