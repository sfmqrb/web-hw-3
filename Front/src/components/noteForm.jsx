import React from "react";
import Joi from "joi-browser";
import Form from "./common/form";
import { getNote, saveNote } from "../services/fakeNoteService";
import { getTypes } from "../services/typeService";
import { toast } from "react-toastify";

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
      // console.log("note", note);
      if (note.data.misscache) {
        toast.warn("MissCached in the GET request happened");
      }
      const localData = this.mapToViewModel(note);
      // console.log("local", localData);
      this.setState({ data: localData });
      // console.log("state", this.state);
    } catch (ex) {
      // console.log("error", ex);
      if (ex.response && ex.response.status === 404)
        return this.props.history.replace("/not-found");
    }
  }

  mapToViewModel(note) {
    return {
      title: note.data.title || "TITLE",
      typeId: note.data.type || "Others",
      text: note.data.text || "TEXT",
      _id: note.data._id || "-1",
    };
  }

  doSubmit = async () => {
    // backend
    // console.log(this.props)
    await saveNote(this.state.data);
    window.location = "/notes";
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
