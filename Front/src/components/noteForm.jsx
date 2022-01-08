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

  componentDidMount() {
    const types = getTypes();
    this.setState({ types });

    const noteId = this.props.match.params.id;
    if (noteId === "new") return;

    const note = getNote(noteId);
    if (!note) return this.props.history.replace("/not-found");

    this.setState({ data: this.mapToViewModel(note) });
  }

  mapToViewModel(note) {
    return {
      _id: note._id,
      title: note.title,
      typeId: note.type._id,
      text: note.text,
    };
  }

  doSubmit = () => {
    saveNote(this.state.data);

    this.props.history.push("/notes");
  };

  render() {
    return (
      <div>
        <h1>Note Form</h1>
        <form onSubmit={this.handleSubmit}>
          {this.renderInput("title", "Title")}
          {this.renderSelect("typeId", "Type", this.state.types)}
          {this.renderInput("text", "Text")}
          {this.renderButton("Save")}
        </form>  
      </div>
    );
  }
}

export default NoteForm;
