import generic from "../generic"

const GetByServiceID = (id) => { return generic.GenericGetAll(`/api/check_service/${id}`) };
const GetByRoundServiceID = (round_id, service_id) => { return generic.GenericGetAll(`/api/check/${round_id}/${service_id}`) };
const GetByRoundID = (id) => { return generic.GenericGetAll(`/api/check_round/${id}`) };



export default {
    GetByServiceID,
    GetByRoundServiceID,
    GetByRoundID
};
