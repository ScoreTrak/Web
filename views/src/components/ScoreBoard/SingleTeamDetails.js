import React, {useEffect, useState, useRef} from "react";
import Accordion from "@material-ui/core/Accordion";
import AccordionSummary from "@material-ui/core/AccordionSummary";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import Typography from "@material-ui/core/Typography";
import AccordionDetails from "@material-ui/core/AccordionDetails";
import Box from "@material-ui/core/Box";
import {makeStyles} from "@material-ui/core/styles";
import CheckCircleOutlineIcon from '@material-ui/icons/CheckCircleOutline';
import ErrorIcon from '@material-ui/icons/Error';
import Alert from "@material-ui/lab/Alert";
import AlertTitle from "@material-ui/lab/AlertTitle";
import Grid from "@material-ui/core/Grid";
import MaterialTable from "material-table";
import PropertyService from "../../services/property/properties";
import CheckService from "../../services/check/check";
import HostService from "../../services/host/hosts";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import Input from "@material-ui/core/Input";
import FormHelperText from "@material-ui/core/FormHelperText";
import Button from "@material-ui/core/Button";
import ConfigService from "../../services/config/config";


const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        width: '100%',

    },

    customAccordionSuccessHeader: {
        borderRight: `1px solid ${theme.palette.success.main}`,
        borderLeft: `1px solid ${theme.palette.success.main}`,
    },
    customAccordionErrorHeader: {
        borderRight: `1px solid ${theme.palette.error.main}`,
        borderLeft: `1px solid ${theme.palette.error.main}`,
    },

    paper: {
        padding: theme.spacing(1),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },

    button: {
        margin: theme.spacing(1),
    },

    input: {
        display: 'none',
    },

    heading: {
        fontSize: theme.typography.pxToRem(15),
        flexBasis: '33.33%',
        flexShrink: 0,
    },
    secondaryHeading: {
        fontSize: theme.typography.pxToRem(15),
        color: theme.palette.text.secondary,
    },
    iconSuccess: {
        color: theme.palette.success.main,
        marginRight: 12,
        fontSize: 22,
        opacity: 0.9,
    },

    iconError: {
        color: theme.palette.error.main,
        marginRight: 12,
        fontSize: 22,
        opacity: 0.9,
    },

}));

export default function SingleTeamDetails(props) {
    const classes = useStyles();
    const [PropertiesData, setPropertiesData] = useState([]);
    const [HostData, setHostData] = useState({});
    const [expanded, setExpanded] = useState(false);
    const [history, setHistory] = useState({});
    const handleChange = (panel) => (event, isExpanded) => {
        setPropertiesData([])
        setHostData({})
        if (!isExpanded){
            setExpanded(false)
        } else {
            setExpanded(panel);
        }
    };
    function usePreviousDT(value) {
        const ref = useRef();
        useEffect(() => {
            ref.current = value;
        });
        return ref.current;
    }
    const prevDT = usePreviousDT({...props.dt});

    useEffect(() => {
        if (prevDT){
            setHistory(prevState => {
                let nextState = {}
                Object.keys(prevState).forEach(cached_service_id => {
                    if (prevState[cached_service_id].length !== 0) {
                        nextState[cached_service_id] = [...prevState[cached_service_id],
                            {
                                service_id: cached_service_id["service_id"],
                                host_id: prevState[cached_service_id][prevState[cached_service_id].length-1]["host_id"],
                                passed: prevDT["Teams"][teamID]["Hosts"][prevState[cached_service_id][prevState[cached_service_id].length-1]["host_id"]]["Services"][cached_service_id]["Passed"],
                                err: prevDT["Teams"][teamID]["Hosts"][prevState[cached_service_id][prevState[cached_service_id].length-1]["host_id"]]["Services"][cached_service_id]["Err"],
                                log: prevDT["Teams"][teamID]["Hosts"][prevState[cached_service_id][prevState[cached_service_id].length-1]["host_id"]]["Services"][cached_service_id]["Log"],
                                round_id: prevDT["Round"],
                            }
                        ]
                    } else {
                        Object.keys(props.dt["Teams"][teamID]["Hosts"]).forEach((host) => {
                            let currentHost = props.dt["Teams"][teamID]["Hosts"][host]
                            Object.keys(currentHost["Services"]).forEach((service_id) => {
                                if (cached_service_id === service_id){
                                    nextState[service_id] = [{
                                        service_id: service_id,
                                        host_id: host,
                                        passed: prevDT["Teams"][teamID]["Hosts"][host]["Services"][service_id]["Passed"],
                                        err: prevDT["Teams"][teamID]["Hosts"][host]["Services"][service_id]["Err"],
                                        log: prevDT["Teams"][teamID]["Hosts"][host]["Services"][service_id]["Log"],
                                        round_id: prevDT["Round"],
                                    }]
                                }
                            })
                        })
                    }
                })
                return {...nextState}
            })
        }
    }, [props.dt]);



    const teamID = props.teamID
    return (
            <Box height="100%" width="100%" align="left" >
                {Object.keys(props.dt["Teams"][teamID]["Hosts"]).map((host) => {
                    let currentHost = props.dt["Teams"][teamID]["Hosts"][host]
                    return Object.keys(currentHost["Services"]).map((service_id) => {
                        let sr = currentHost["Services"][service_id]
                        let keyName
                        if (sr["DisplayName"]){
                            keyName = sr["DisplayName"]
                        } else {
                            if (currentHost["HostGroup"]){
                                keyName =currentHost["HostGroup"]["Name"] + "-" + sr["Name"]
                            } else{
                                keyName = sr["Name"]
                            }
                        }
                        return (
                            <Accordion expanded={expanded === keyName} onChange={handleChange(keyName)} className={sr["Passed"] ? classes.customAccordionSuccessHeader: classes.customAccordionErrorHeader}>
                                <AccordionSummary
                                    expandIcon={<ExpandMoreIcon />}
                                    aria-controls={`${keyName}bh-content`}
                                    id={`${keyName}bh-header`}
                                >
                                    {sr["Passed"] ? <CheckCircleOutlineIcon className={classes.iconSuccess}  />  : <ErrorIcon className={classes.iconError}/>}
                                    <Typography className=  {classes.heading}>{keyName}</Typography>
                                    <Typography className={classes.secondaryHeading}>Host used for last round: {currentHost["Address"]}</Typography>

                                </AccordionSummary>
                                <AccordionDetails>
                                    {expanded === keyName &&
                                        <SingleTeamDetailsAccordionDetailsBox {...props} history={history} prevDT={prevDT} setHistory={setHistory} teamID={teamID} host_id={host} service_id={service_id} sr={sr} setHostData={setHostData} HostData={HostData}  PropertiesData={PropertiesData} setPropertiesData={setPropertiesData} />
                                    }
                                </AccordionDetails>
                            </Accordion>
                        );
                    })
                })}
            </Box>
    );
}


function SingleTeamDetailsAccordionDetailsBox(props) {
    const sr = props.sr
    const PropertiesData = props.PropertiesData
    const classes = useStyles();
    const setPropertiesData = props.setPropertiesData
    const service_id = props.service_id

    const history = props.history
    const setHistory = props.setHistory

    const host_id = props.host_id

    const columns = [
        { title: 'Key', field: 'key', editable: "never"},
        { title: 'Value Used', field: 'value_used', editable:"never"},
        { title: 'Current Value', field: 'value' },
    ]

    const columnsPreviousRounds = [
        { render: rowData => <div> {rowData.passed ? <CheckCircleOutlineIcon className={classes.iconSuccess}  />  : <ErrorIcon className={classes.iconError}/>}  </div> },
        { title: 'Round', field: 'round_id', defaultSort: "desc"},
        { title: 'Passed', field: 'passed', hidden: true},
        { title: 'Parent Host ID', field: 'host_id', hidden: true},
        { title: 'Service ID', field: 'service_id', hidden: true},
        { title: 'Response', field: 'log' },
        { title: 'Error Details', field: 'err' },
    ]

    async function reloadPreviousChecks(service) {
        const results = await CheckService.GetByServiceID(service)
        if ((!!results) && (results.constructor === Object)){
            return []
        }
        return [...results]
    }


    async function reloadProperties(service) {
        const results = await PropertyService.GetAllByServiceID(service)
        if ((!!results) && (results.constructor === Object)){
            return []
        }
        return [...results]
    }

    async function reloadHost(hostID) {
        const results = await HostService.GetByID(hostID)
        return {...results}
    }

    function reloadPropertiesSetter(service_id, sr) {
        reloadProperties(service_id).then( results => {
            let d = []
            for (const [key, property] of Object.entries(sr["Properties"])) {
                let obj = {key: key, value_used: property["Value"], service_id: key}
                results.forEach(res => {
                    if (key === res["key"] && res["status"] === "Edit"){
                        obj.value = res["value"]
                    }
                })
                d.push(obj)
            }
            setPropertiesData(d)
        })
    }

    const handleSetHostAddress = (e, hstID) => {
        e.preventDefault()
        let val = document.getElementById(`host_address_${hstID}`).value
        props.setHostData(prevState => {return{...prevState, address: val}})
        HostService.Update(host_id, {address: val}).then( async () => {
            reloadHostSetter(hstID)
        }, props.errorSetter)
    }

    function reloadHostSetter(host_id) {
        reloadHost(host_id).then( results => {
            props.setHostData(results)
        })
    }

    useEffect(() => {
        if (!history[service_id]){
            reloadPreviousChecks(service_id).then( results => {
                let d = []
                results.forEach(res => {
                    if (res["round_id"] < props.dt.Round){
                        d.push({service_id: service_id, round_id: res["round_id"], passed: res["passed"], log: res["log"], err: res["err"], host_id: host_id})
                    }
                })
                setHistory(prevState => {return {...prevState, [service_id]: d }})
            })
        }
        reloadHostSetter(host_id)
        reloadPropertiesSetter(service_id, sr)
    }, []);

    return (
        <Box width="100%" bgcolor="background.paper">
            <Grid container spacing={3}>
                <Grid item xs={6}>
                    {sr["Log"] &&
                    <Alert severity={sr["Passed"] ? "info" : "warning"}>
                        <AlertTitle>Response</AlertTitle>
                        {sr["Log"]}
                    </Alert>
                    }
                    <br/>
                    {sr["Err"] &&
                    <Alert severity="error">
                        <AlertTitle>Error Details</AlertTitle>
                        {sr["Err"]}
                    </Alert>
                    }
                </Grid>
                <Grid item xs={6}>
                    <MaterialTable
                        options={{pageSizeOptions: [3, 5,10,20,50,100], pageSize:3}}
                        title="Properties"
                        columns={columns}
                        data={PropertiesData}
                        cellEditable={{
                            onCellEditApproved: (newValue, oldValue, rowData, columnDef) => {
                                return new Promise((resolve) => {
                                    setTimeout(() => {
                                        resolve();
                                        setPropertiesData((prevState) => {
                                            return prevState.map(property => {
                                                if (property["key"] === rowData["key"]) {
                                                    return {...rowData, value: newValue}
                                                }
                                                return {...property}
                                            })
                                        });
                                        PropertyService.Update(service_id, rowData["key"],{'value': newValue}).then( async () =>{
                                            reloadPropertiesSetter(service_id, sr)
                                        }, (error) => {
                                            props.errorSetter(error)
                                        })
                                    }, 600);
                                })
                            }
                        }}
                    />
                    {
                        props.HostData["edit_host"] &&
                        <form style={{width: "100%", marginTop: "1vh"}} onSubmit={e => {handleSetHostAddress(e, host_id) }}>
                            <FormControl style={{ display: 'flex', flexDirection: 'row', width: "100%"}}>
                                <div>
                                    <InputLabel htmlFor="host_address">Host (Current: {props.HostData["address"]})</InputLabel>
                                    <Input id={`host_address_${host_id}`} aria-describedby="my-helper-text" />
                                    <FormHelperText id="my-helper-text">Set the address of the remote machine</FormHelperText>
                                </div>
                                <Button type="submit" variant="outlined" color="primary" style={{width: "10vh", height: "3vh", marginLeft: "3vh", marginTop: "auto", marginBottom:"auto"}}>
                                    Set
                                </Button >
                            </FormControl>
                        </form>
                    }
                </Grid>
            </Grid>
            <Grid container spacing={3}>
                <Grid item xs={12}>
                    <MaterialTable
                        options={{pageSizeOptions: [5,10,20,50,100], pageSize:5}}
                        title="Previous Rounds"
                        columns={columnsPreviousRounds}
                        data={history[service_id]}
                    />
                </Grid>
            </Grid>
        </Box>
    )
}

//Todo: Allow changing hostname if editable.