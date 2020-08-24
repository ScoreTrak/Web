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


    const [data, setData] = useState([]);
    const [expanded, setExpanded] = React.useState(false);

    const handleChange = (panel, service_id, service) => (event, isExpanded) => {
        setData([])
        if (!isExpanded){
            setExpanded(false)
        } else {
            setExpanded(panel);
        }
    };

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
                            <Accordion expanded={expanded === keyName} onChange={handleChange(keyName, service_id, sr)} className={sr["Passed"] ? classes.customAccordionSuccessHeader: classes.customAccordionErrorHeader}>
                                <AccordionSummary
                                    expandIcon={<ExpandMoreIcon />}
                                    aria-controls={`${keyName}bh-content`}
                                    id={`${keyName}bh-header`}
                                >
                                    {sr["Passed"] ? <CheckCircleOutlineIcon className={classes.iconSuccess}  />  : <ErrorIcon className={classes.iconError}/>}
                                    <Typography className=  {classes.heading}>{keyName}</Typography>
                                    <Typography className={classes.secondaryHeading}>Host: {currentHost["Address"]}</Typography>

                                </AccordionSummary>
                                <AccordionDetails>
                                    {expanded === keyName &&
                                        <SingleTeamDetailsAccordionDetailsBox {...props} teamID={teamID} host_id={host} service_id={service_id} sr={sr} data={data} setData={setData} />
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
    const data = props.data
    const classes = useStyles();
    const setData = props.setData

    const service_id = props.service_id
    const teamID = props.teamID
    const host_id = props.host_id


    const [history, setHistory] = React.useState([]);



    const columns = [
        { title: 'Key', field: 'key', editable: "never"},
        { title: 'Value Used', field: 'value_used', editable:"never"},
        { title: 'Current Value', field: 'value' },
    ]

    const columnsPreviousRounds = [
        { title: 'Passed Image', render: rowData => <div> {rowData.passed ? <CheckCircleOutlineIcon className={classes.iconSuccess}  />  : <ErrorIcon className={classes.iconError}/>}  </div> },
        { title: 'Round', field: 'round_id', defaultSort: "desc"},
        { title: 'Passed', field: 'passed', hidden: true},
        { title: 'Service ID', field: 'service_id', hidden: true},
        { title: 'Response', field: 'log' },
        { title: 'Error Details', field: 'err' },
    ]

    const prevDT = usePreviousDT({...props.dt});

    async function reloadPreviousChecks(service) {
        const results = await CheckService.GetByServiceID(service)
        if ((!!results) && (results.constructor === Object)){
            return []
        }
        return [...results]
    }

    function usePreviousDT(value) {
        const ref = useRef();
        useEffect(() => {
            ref.current = value;
        });
        return ref.current;
    }


    async function reloadProperties(service) {
        const results = await PropertyService.GetAllByServiceID(service)
        if ((!!results) && (results.constructor === Object)){
            return []
        }
        return [...results]
    }

    function reloadSetter(service_id, sr) {
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
            setData(d)
        })
    }

    useEffect(() => {
        reloadPreviousChecks(service_id).then( results => {
            let d = []
            results.forEach(res => {
                if (res["round_id"] !== props.dt.Round){
                    d.push({service_id: service_id, round_id: res["round_id"], passed: res["passed"], log: res["log"], err: res["err"]})
                }
            })
            setHistory(d)
        })
        reloadSetter(service_id, sr)
        const roundReload = setInterval(() => {
            reloadSetter(service_id, sr)
        }, 10000);
        return () => clearInterval(roundReload);

    }, []);

    useEffect(() => {
        if (prevDT){
            console.log(prevDT)
            console.log(props.dt)
            console.log(history)
            console.log("Setting History inside USEFEECT SECONDARY")
            setHistory(prevState => {return [...prevState,
                {
                    service_id: service_id["service_id"],
                    passed: prevDT["Teams"][teamID]["Hosts"][host_id]["Services"][service_id]["Passed"],
                    err: prevDT["Teams"][teamID]["Hosts"][host_id]["Services"][service_id]["Err"],
                    log: prevDT["Teams"][teamID]["Hosts"][host_id]["Services"][service_id]["Log"],
                    round_id: prevDT["Round"],
                }
            ]})
        }
    }, [props.dt]);

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
                        options={{pageSizeOptions: [5,10,20,50,100], pageSize:20, emptyRowsWhenPaging:false}}
                        title="Properties"
                        columns={columns}
                        data={data}
                        cellEditable={{
                            onCellEditApproved: (newValue, oldValue, rowData, columnDef) => {
                                return new Promise((resolve) => {
                                    setTimeout(() => {
                                        resolve();
                                        setData((prevState) => {
                                            return prevState.map(property => {
                                                if (property["key"] === rowData["key"]) {
                                                    return {...rowData, value: newValue}
                                                }
                                                return {...property}
                                            })
                                        });
                                        PropertyService.Update(service_id, rowData["key"],{'value': newValue}).then( async () =>{
                                            reloadSetter(service_id, sr)
                                        }, (error) => {
                                            props.errorSetter(error)
                                        })
                                    }, 600);
                                })
                            }
                        }}
                    />

                </Grid>
            </Grid>
            <Grid container spacing={3}>
                <Grid item xs={12}>
                    <MaterialTable
                        options={{pageSizeOptions: [5,10,20,50,100], pageSize:5}}
                        title="Previous Rounds"
                        columns={columnsPreviousRounds}
                        data={history}
                    />
                </Grid>
            </Grid>
        </Box>
    )

}