import React from "react";
import Joi from "joi-browser";
import Form from "./common/form";
import {login} from "../services/authService";

class LoginForm extends Form {
    state = {
        data: {username: "", password: ""},
        errors: {},
    };

    schema = {
        username: Joi.string().required().label("Username"),
        password: Joi.string().required().label("Password"),
    };

    doSubmit = async () => {
        // backend
        console.log("login");
        const {data} = this.state;
        console.log(data);
        const output = await login(data.username, data.password);
        console.log(output);
        localStorage.setItem("jwt", output.data["jwt"]);
        localStorage.setItem("notes", JSON.stringify(output.data["notes"]));
        // //debug
        // const notes = output.data["notes"]
        // console.log(notes)
        // console.log(notes[0])
        // const notesLocal = localStorage.getItem("notes")
        // console.log("local")
        // console.log(notesLocal)
        // console.log(notesLocal[0])
        window.location = "/"; // full reload of app
        // try catch 400 username not exists
        console.log("Submitted");
    };

    render() {
        return (
            <div>
                <h1>Login</h1>
                <form onSubmit={this.handleSubmit}>
                    {this.renderInput("username", "Username")}
                    {this.renderInput("password", "Password", "password")}
                    {this.renderButton("Login")}
                </form>
            </div>
        );
    }
}

export default LoginForm;
