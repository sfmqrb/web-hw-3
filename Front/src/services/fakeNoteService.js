import * as typesAPI from "./typeService";

const notes = [
  {
    _id: "5b21ca3eeb7f6fbccd471815",
    title: "Terminator",
    type: { _id: "Personal", name: "Personal" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd471816",
    title: "Die Hard",
    type: { _id: "Personal", name: "Personal" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd471817",
    title: "Get Out",
    type: { _id: "Work", name: "Work" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd471819",
    title: "Trip to Italy",
    type: { _id: "Family", name: "Family" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181a",
    title: "Airplane",
    type: { _id: "Family", name: "Family" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181b",
    title: "Wedding Crashers",
    type: { _id: "Family", name: "Family" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181e",
    title: "Gone Girl",
    type: { _id: "Work", name: "Work" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd47181f",
    title: "The Sixth Sense",
    type: { _id: "Work", name: "Work" },
    text: "This is a test",
  },
  {
    _id: "5b21ca3eeb7f6fbccd471821",
    title: "The Avengers",
    type: { _id: "Personal", name: "Personal" },
    text: "This is a test",
  },
];

export function getNotes() {
  return notes;
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
  return noteInDb;
}
