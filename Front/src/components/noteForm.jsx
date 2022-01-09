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
  componentDidMount() {
    const types = getTypes();
    this.setState({ types });

    const noteId = this.props.match.params.id;
    if (noteId === "new") return;

    // client
    // const note = getNote(noteId);
    // backend
    try {
      const note = await getNote(noteId);
      this.setState({ data: this.mapToViewModel(note) });
    } catch (ex) {
      if (ex.response && ex.response.status === 404)
        return this.props.history.replace("/not-found");
    }

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
