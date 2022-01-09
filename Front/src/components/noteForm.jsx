import React from "react";
import Joi from "joi-browser";
import Form from "./common/form";
import { getNote, saveNote } from "../services/fakeNoteService";
import { getTypes } from "../services/typeService";

class NoteForm extends Form {
  state = {
    data: {
      title: "",
      typeId: "",
      text: "",
    },
    types: [],
    errors: {},
  };

  schema = {
    _id: Joi.string(),
    title: Joi.string().required().label("Title"),
    typeId: Joi.string().required().label("Type"),
    text: Joi.string().required().label("Text"),
  };

  // async //backend
  async componentDidMount() {
    const types = getTypes();
    this.setState({ types });

    const noteId = this.props.match.params.id;
    if (noteId === "new") return;

    // client
    // const note = getNote(noteId);
    // backend
    try {
      const note = await getNote(noteId);

      console.log("get note main function try", note);
      const localData = this.mapToViewModel(note);
      this.setState({ data: localData });
    } catch (ex) {
      console.log("get note main function catch");
      if (ex.response && ex.response.status === 404)
        return this.props.history.replace("/not-found");
    }
  }

  mapToViewModel(note) {
    console.log("mapToViewModel", note);
    console.log(note.data);
    console.log(note.data.title || "title");

    return {
      title: note.data.title || "TITLE",
      typeId: "Others",
      text: note.data.text || "TEXT",
      _id: Number(note.data._id) || -1,

      // _id: 1,
      // title: "title",
      // typeId: "Others",
      // text: "text",
    };
  }

  doSubmit = () => {
    saveNote(this.state.data);
    // backend
    // await saveNote(this.state.data);

    this.props.history.push("/notes");
    // backend
  };

  render() {
    return (
      <div>
        <h1>Note Form</h1>
        <form onSubmit={this.handleSubmit}>
          {this.renderInput("title", "Title")}
          {this.renderSelect("typeId", "Type", this.state.types)}
          {this.renderInput("text", "Text", "text", "textarea")}
          {this.renderButton("Save")}
        </form>
      </div>
    );
  }
}

export default NoteForm;
