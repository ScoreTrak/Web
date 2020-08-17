import React, {useEffect, useState} from "react";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import Box from '@material-ui/core/Box';
import ConfigService from "../../services/config/config";
import CompetitionService from "../../services/competition/competition";
import PolicyService from "../../services/policy/policy";
import ReactJson from 'react-json-view'
import CircularProgress from "@material-ui/core/CircularProgress";
import { makeStyles } from '@material-ui/core/styles';
import SaveIcon from '@material-ui/icons/Save';
import Button from '@material-ui/core/Button';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import Slide from "@material-ui/core/Slide";
import Dialog from "@material-ui/core/Dialog";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogActions from "@material-ui/core/DialogActions";
import Accordion from '@material-ui/core/Accordion';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Switch from "@material-ui/core/Switch";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import Input from "@material-ui/core/Input";
import FormHelperText from "@material-ui/core/FormHelperText";

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        width: '100%',
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

}));


const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});


export default function Settings(props) {
    const classes = useStyles();
    const setTitle = props.setTitle
    const classesPaper = props.classesPaper
    setTitle("Settings")
    const [dt, setData] = useState({
        loader:true, config: undefined,
        staticConfig: undefined, staticWebConfig: undefined,
        policy: undefined
    });



    const [open, setOpen] = React.useState("");
    const [fileSelected, setFileSelected] = React.useState({selected: false, name: ""});

    const handleClickOpen = (panel) => () => {
        setOpen(panel);
    };

    const handleClose = () => {
        setOpen("");
    };

    const [expanded, setExpanded] = React.useState('panelConfig');

    const handleChange = (panel) => (event, isExpanded) => {
        setExpanded(isExpanded ? panel : false);
    };
    
    
    const handleSetFileSelected = () => {
        setFileSelected({selected: true, name: document.getElementById('file').files[0].name})
    }


    async function reload() {
        const config = await ConfigService.Get()
        const policy = await PolicyService.Get()
        return {config: config, policy: policy}
    }

    async function loadAll(){
        const staticConfig = await ConfigService.GetStaticConfig()
        const staticWebConfig = await ConfigService.GetStaticWebConfig()
        return {...(await reload()), staticConfig: staticConfig, staticWebConfig: staticWebConfig}
    }
    useEffect(() => {
        loadAll().then(newState => { setData(prevState => {return{...prevState, ...newState, loader:false}})}, props.errorSetter)
    }, []);

    const handleSetEnabled = (e) => {
        setData(prevState => {return{...prevState, config:{...prevState.config, enabled: e.target.checked}}})
        ConfigService.Update({enabled: e.target.checked}).then( async () => {
                reload().then(newState => { setData(prevState => {return{...prevState, ...newState}})}, props.errorSetter)
        }, props.errorSetter)
    }

    const handleSetRoundDuration = (e) => {
        e.preventDefault()
        let val = Number(document.getElementById("round_duration").value)
        setData(prevState => {return{...prevState, config:{...prevState.config, round_duration: val}}})
        ConfigService.Update({round_duration: val}).then( async () => {
            reload().then(newState => { setData(prevState => {return{...prevState, ...newState}})}, props.errorSetter)
        }, props.errorSetter)

    }


    const handleSetPolicy = (e) => {
        setData(prevState => {return{...prevState, policy:{...prevState.policy, [e.target.value]: e.target.checked}}})
        PolicyService.Update({[e.target.value]: e.target.checked}).then( async () => {
            reload().then(newState => { setData(prevState => {return{...prevState, ...newState}})}, props.errorSetter)
        }, props.errorSetter)
    }

    const handleUpload = () => {
        let formData = new FormData();
        formData.append("file", document.getElementById('file').files[0]);
        CompetitionService.LoadCompetition(formData).then(() => {
            document.getElementById('file').value = ""
            loadAll().then(newState => { setData(prevState => {return{...prevState, ...newState}})}, props.errorSetter)
            setFileSelected({selected: false, name: ""})
        }, props.errorSetter)
        handleClose()
    }

    const handleResetCompetition = () => {
        ConfigService.ResetCompetition().then(() => {
            loadAll().then(newState => { setData(prevState => {return{...prevState, ...newState}})}, props.errorSetter)
            setFileSelected({selected: false, name: ""})
        }, props.errorSetter)
        handleClose()
    }

    const handleDeleteCompetition = () => {
        ConfigService.DeleteCompetition().then(() => {
            loadAll().then(newState => { setData(prevState => {return{...prevState, ...newState}})}, props.errorSetter)
            setFileSelected({selected: false, name: ""})
        }, props.errorSetter)
        handleClose()
    }

    return (
        <Paper className={classesPaper} style={{minHeight: "85vh"}}>
            {!dt.loader ?
                <Box height="100%" width="100%" align="left" >
                    <Accordion expanded={expanded === 'panelConfig'} onChange={handleChange('panelConfig')}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panelConfigbh-content"
                            id="panelConfigbh-header"
                        >
                            <Typography className={classes.heading}>ScoreTrak Settings</Typography>
                            <Typography className={classes.secondaryHeading}>Following are the dynamically configurable settings for scoring masters</Typography>

                        </AccordionSummary>
                        <AccordionDetails>
                            <Box p={1} m={1} bgcolor="background.paper">

                            <FormControlLabel
                                control={
                            <Switch checked={dt.config["enabled"]} onChange={handleSetEnabled} />
                                }
                                label="Enable Competition?"
                            />

                            <br/>

                                <form onSubmit={handleSetRoundDuration}>
                                    <FormControl>
                                        <InputLabel htmlFor="round_duration">Round Duration (Current: {dt.config["round_duration"]})</InputLabel>
                                        <Input id="round_duration" aria-describedby="my-helper-text" type="number" inputProps={{ min: "20"}} />
                                        <FormHelperText id="my-helper-text">Number of seconds it takes for one round to elapse.</FormHelperText>
                                        <Button type="submit" variant="outlined" color="primary" >
                                            Set
                                        </Button>
                                    </FormControl>
                                </form>

                            </Box>
                        </AccordionDetails>
                    </Accordion>

                    <Accordion expanded={expanded === 'panelPolicy'} onChange={handleChange('panelPolicy')}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panelPolicybh-content"
                            id="panelPolicybh-header"
                        >
                            <Typography className={classes.heading}>Policy</Typography>
                            <Typography className={classes.secondaryHeading}>Configure policies for allowing/disallowing resources</Typography>

                        </AccordionSummary>
                        <AccordionDetails>
                            <Box p={1} m={1} bgcolor="background.paper">
                            <FormControlLabel
                                control={
                                    <Switch checked={dt.policy["allow_unauthenticated_users"]} onChange={handleSetPolicy} value="allow_unauthenticated_users" />
                                }
                                label="Allow unauthenticated users to see scoreboard"
                            />
                            <br />
                            <FormControlLabel
                                control={
                                    <Switch checked={dt.policy["allow_changing_usernames_and_passwords"]} onChange={handleSetPolicy} value="allow_changing_usernames_and_passwords" />
                                }
                                label="Allow users to change usernames and passwords"
                            />
                                <br />
                            <FormControlLabel
                                control={
                                    <Switch checked={dt.policy["allow_to_see_points"]} onChange={handleSetPolicy} value="allow_to_see_points" />
                                }
                                label="Allow users to see other teams' points"
                            />
                                <br />
                            <FormControlLabel
                                control={
                                    <Switch checked={dt.policy["show_addresses"]} onChange={handleSetPolicy} value="show_addresses" />
                                }
                                label="Allow users to see other teams' addresses"
                            />
                            </Box>

                        </AccordionDetails>
                    </Accordion>
                    
                    <Accordion expanded={expanded === 'panelStaticConfig'} onChange={handleChange('panelStaticConfig')}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panelStaticConfigbh-content"
                            id="panelStaticConfigbh-header"
                        >
                            <Typography className={classes.heading}>Static ScoreTrak Config</Typography>
                            <Typography className={classes.secondaryHeading}>
                                This is a JSON representation of the Static Config for Scoring Master
                            </Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Box component="span" display="block" p={1} m={1} bgcolor="background.paper">
                                <ReactJson src={dt.staticConfig} style={{backgroundColor: "inherit"}} onDelete={false} onEdit={false} displayDataTypes={false} displayObjectSize={false} theme={props.isDarkTheme ? "monokai" : "bright:inverted"}/>
                            </Box>
                        </AccordionDetails>
                    </Accordion>
                    <Accordion expanded={expanded === 'panelStaticWebConfig'} onChange={handleChange('panelStaticWebConfig')}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panelStaticWebConfigbh-content"
                            id="panelStaticWebConfigbh-header"
                        >
                            <Typography className={classes.heading}>Static Web Config</Typography>
                            <Typography className={classes.secondaryHeading}>
                                This is a JSON representation of the Static Config for Web component
                            </Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Box component="span" display="block" p={1} m={1} bgcolor="background.paper">
                                <ReactJson src={dt.staticWebConfig} style={{backgroundColor: "inherit"}} onDelete={false} onEdit={false} displayDataTypes={false} displayObjectSize={false} theme={props.isDarkTheme ? "monokai" : "bright:inverted"}/>
                            </Box>
                        </AccordionDetails>
                    </Accordion>
                    <Accordion expanded={expanded === 'panelImportExport'} onChange={handleChange('panelImportExport')}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panelImportExportbh-content"
                            id="panelImportExportbh-header"
                        >
                            <Typography className={classes.heading}>Export/Import Competition</Typography>
                            <Typography className={classes.secondaryHeading}>
                                Export or Import the competition using a JSON representation of the competition
                            </Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Box align="center" width="100%" p={1} m={1} bgcolor="background.paper">
                                <Button
                                    variant="contained"
                                    color="primary"
                                    className={classes.button}
                                    startIcon={<SaveIcon />}
                                    onClick={CompetitionService.FetchCoreCompetition}
                                >
                                    Export Core Competition
                                </Button>
                                <Button variant="outlined" color="primary" onClick={handleClickOpen("upload")}>
                                    Upload Competition
                                </Button>
                                <Dialog
                                    open={open === 'upload'}
                                    TransitionComponent={Transition}
                                    keepMounted
                                    onClose={handleClose}
                                    aria-labelledby="alert-dialog-slide-title"
                                    aria-describedby="alert-dialog-slide-description"
                                >
                                    <DialogContent>
                                        <DialogContentText id="alert-dialog-slide-description" align="center">
                                            <input
                                                className={classes.input}
                                                id="file"
                                                type="file"
                                                onChange={handleSetFileSelected}
                                            />
                                            <label htmlFor="file">
                                                <Button
                                                    variant="contained"
                                                    color="primary"
                                                    className={classes.button}
                                                    component="span"
                                                    startIcon={<CloudUploadIcon />}
                                                >
                                                    Load Competition
                                                </Button>
                                            </label>
                                        </DialogContentText>
                                        { fileSelected.selected &&
                                        <Typography align="center">{fileSelected.name}</Typography>
                                        }
                                    </DialogContent>
                                    {fileSelected &&
                                    <DialogActions>
                                        <Button onClick={handleClose} color="primary">
                                            Cancel
                                        </Button>
                                        <Button onClick={handleUpload} color="primary">
                                            Upload
                                        </Button>
                                    </DialogActions>
                                    }
                                </Dialog>
                                <Button
                                    variant="contained"
                                    color="primary"
                                    className={classes.button}
                                    startIcon={<SaveIcon />}
                                    onClick={CompetitionService.FetchEntireCompetition}
                                >
                                    Export Entire Competition
                                </Button>
                            </Box>
                        </AccordionDetails>
                    </Accordion>

                    <Accordion expanded={expanded === 'panelDeleteReset'} onChange={handleChange('panelDeleteReset')}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="panelDeleteReseth-content"
                            id="panelDeleteResetbh-header"
                        >
                            <Typography className={classes.heading}>Reset/Delete Competition</Typography>
                            <Typography className={classes.secondaryHeading}>
                                Reset: Resets Scores, and Rounds. Delete: Removes everything
                            </Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            <Box align="center" width="100%" p={1} m={1} bgcolor="background.paper">
                                <Button variant="outlined" style={{color: "red", border: '1px solid red' }} onClick={handleClickOpen("reset")} className={classes.button}>
                                    Reset Competition
                                </Button>
                                <Dialog
                                    open={open === "reset"}
                                    TransitionComponent={Transition}
                                    keepMounted
                                    onClose={handleClose}
                                    aria-labelledby="alert-dialog-slide-title"
                                    aria-describedby="alert-dialog-slide-description"
                                >
                                    <DialogContent>
                                        <DialogContentText id="alert-dialog-slide-description" align="center">
                                            Are you sure?
                                        </DialogContentText>
                                    </DialogContent>

                                    <DialogActions>
                                        <Button onClick={handleClose} color="primary">
                                            Cancel
                                        </Button>
                                        <Button onClick={handleResetCompetition} color="primary">
                                            Reset Competition
                                        </Button>
                                    </DialogActions>

                                </Dialog>


                                <Button variant="outlined" style={{color: "red", border: '1px solid red' }} onClick={handleClickOpen("delete")} className={classes.button}>
                                    Delete Competition
                                </Button>
                                <Dialog
                                    open={open === "delete"}
                                    TransitionComponent={Transition}
                                    keepMounted
                                    onClose={handleClose}
                                    aria-labelledby="alert-dialog-slide-title"
                                    aria-describedby="alert-dialog-slide-description"
                                >
                                    <DialogContent>
                                        <DialogContentText id="alert-dialog-slide-description" align="center">
                                           Are you sure?
                                        </DialogContentText>
                                    </DialogContent>

                                    <DialogActions>
                                        <Button onClick={handleClose} color="primary">
                                            Cancel
                                        </Button>
                                        <Button onClick={handleDeleteCompetition} color="primary">
                                            Delete Competition
                                        </Button>
                                    </DialogActions>

                                </Dialog>
                            </Box>
                        </AccordionDetails>
                    </Accordion>

                </Box>
                :
                <Box height="100%" width="100%" m="auto">
                    <CircularProgress  />
                </Box>
            }

        </Paper>
    );
}