import * as typesAPI from "./typeService";

const notes = [
  {
    _id: "5b21ca3eeb7f6fbccd471815",
    title: "Terminator",
    // type: { _id: "Personal", name: "Personal" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
    misscached: false,
  },
  {
    _id: "5b21ca3eeb7f6fbccd471816",
    title: "Die Hard",
    // type: { _id: "Personal", name: "Personal" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
  {
    _id: "5b21ca3eeb7f6fbccd471817",
    title: "Get Out",
    // type: { _id: "Work", name: "Work" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
  {
    _id: "5b21ca3eeb7f6fbccd471819",
    title: "Trip to Italy",
    // type: { _id: "Family", name: "Family" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181a",
    title: "Airplane",
    type: { _id: "Family", name: "Family" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181b",
    title: "Wedding Crashers",
    // type: { _id: "Family", name: "Family" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181e",
    title: "Gone Girl",
    // type: { _id: "Work", name: "Work" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181f",
    title: "The Sixth Sense",
    // type: { _id: "Work", name: "Work" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
  {
    _id: "5b21ca3eeb7f6fbccd471821",
    title: "The Avengers",
    // type: { _id: "Personal", name: "Personal" },
    text: "lapsus (English) Origin & history From Latin lapsus. Pronunciation IPA: /ˈlæpsəs/ Noun lapsus (pl. lapsus) A slip, lapse, or error. Derived words & phrases lapsus… delapsum: delapsum (Latin) Participle dēlapsum Inflection of dēlapsus (nominative neuter singular) Inflection of dēlapsus (accusative masculine singular) Inflection of dēlapsus (accusative",
  },
];

export function getNotes() {
  // if (!localStorage.getItem("jwt")) {
  //   console.log("first login or register");
  //   window.location = "/login";
  //   // this.props.history.push("/login");
  //   // console.log()
  //   return [];
  // }
  // if logged in ok else register or logged in
  // return {};
  localStorage.setItem("notes", JSON.stringify(notes));
  return localStorage.getItem("notes");
}

export function getNote(id) {
  return notes.find((m) => m._id === id);
}

export function saveNote(note) {
  let noteInDb = notes.find((m) => m._id === note._id) || {};
  noteInDb.title = note.title;
  noteInDb.type = typesAPI.types.find((g) => g._id === note.typeId);
  noteInDb.text = note.text;

  if (!noteInDb._id) {
    noteInDb._id = Date.now().toString();
    notes.push(noteInDb);
  }

  return noteInDb;
}

export function deleteNote(id) {
  let noteInDb = notes.find((m) => m._id === id);
  notes.splice(notes.indexOf(noteInDb), 1);
  // backend
  return noteInDb;
}
