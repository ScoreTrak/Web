import generic from "../generic"

const API_URL = "/api/team/";


const Create = (teams) => { return generic.GenericCreate(API_URL, teams) };
const Update = (id, team) => { return generic.GenericUpdate(API_URL, id, team) };
const GetByID = (id) => { return generic.GenericGetByID(API_URL, id) };
const GetAll = () => { return generic.GenericGetAll(API_URL) };
const Delete = (id) => { return generic.GenericDelete(API_URL, id) };


export default {
    Create,
    Update,
    GetByID,
    GetAll,
    Delete
};
