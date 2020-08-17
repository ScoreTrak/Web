import generic from "../generic";

export const GetLastNoneElapsingRound = () => {
    return generic.GenericGetAll("/api/last_non_elapsing/")
};

export const GetAll = () => {
    return generic.GenericGetAll("/api/round/")
};

export default {
    GetLastNoneElapsingRound,
    GetAll
}