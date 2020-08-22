import generic from "../generic"

const API_URL = "/api/service/";


const Create = (service) => { return generic.GenericCreate(API_URL, service) };
const Update = (id, service) => { return generic.GenericUpdate(API_URL, id, service) };
const GetByID = (id) => { return generic.GenericGetByID(API_URL, id) };
const TestService = (id) => { return generic.GenericGetByID("/api/service_test/", id) };
const GetAll = () => { return generic.GenericGetAll(API_URL) };
const Delete = (id) => { return generic.GenericDelete(API_URL, id) };


export default {
    Create,
    Update,
    GetByID,
    GetAll,
    Delete,
    TestService
};
