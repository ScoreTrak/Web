import React, {useState, useRef, useEffect} from "react";
import ReportService from "../../services/report/report"
import RoundService from "../../services/round/round"
import CircularProgress from '@material-ui/core/CircularProgress';
import Box from '@material-ui/core/Box';
import {Typography} from "@material-ui/core";
import { Route } from "react-router-dom";
import Ranks from "./Ranks";
import Status from "./Status";
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import Button from '@material-ui/core/Button';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import {makeStyles} from "@material-ui/core/styles";
import Paper from "@material-ui/core/Paper";
import AuthService from "../../services/auth/auth";
import Details from "./Details";

const fullScreenButtonStyle = {
    position: 'absolute',
    maxHeight: '3vh',
    bottom: '1vh',
    marginRight: '2vh',
    right: '20px',
}

const useStyles = makeStyles((theme) => ({
    "hidden": { opacity: 0.1, transition: "opacity 0.2s linear"},
    "visible": { opacity: 1, transition: "opacity 0.2s linear"},
    alignItemsAndJustifyContent: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
}));

export default function ScoreBoard(props) {
    const [dt, setData] = useState({Round: 0, loader:true, notStarted: false});
    const [visible, setFade] = React.useState(false);
    const classes = useStyles();
    const handleFadeOver = () => {
        setFade(true);
    };
    const handleFadeOut = () => {
        setFade(false);
    };


    const roundRef = useRef(0)
    let darkTheme = props.isDarkTheme
    const handleFullScreen = props.handleFullScreen
    useEffect(() => {
        function reload() {
            RoundService.GetLastNoneElapsingRound().then( round =>{
                if (round.id !== roundRef.current){
                    ReportService.Get().then( simpleReport => {
                        roundRef.current = simpleReport.Round
                        setData(simpleReport)
                        props.setTitle("Round " + simpleReport.Round)
                    })
                }

            },
                (error) => {
                    if (error.response.status === 404){
                        setData({loader: true, Round: 0, notStarted: true})
                    } else{
                        props.errorSetter(error)
                        setData({loader: true, Round: 0, notStarted: false})
                    }
                })
        }
        reload()
        const roundReload = setInterval(reload, 3000);
        return () => clearInterval(roundReload);
    }, []);

    return (
        <Paper className={handleFullScreen.active ? props.fullSizeHeightPaper : props.fixedHeightPaper}>
            <Box className={classes.alignItemsAndJustifyContent} height="100%" width="100%"  >
                { !("loader" in dt) && dt.Round !==0 ?
                    <Box m="auto" height="100%" width="100%">
                        {
                            (props.currentPolicy["allow_to_see_points"] || AuthService.getCurrentRole() === "black") &&
                        <Route exact path='/ranks' render={() => (
                            <Ranks isDarkTheme={darkTheme} dt={dt}/>
                        )}/>
                        }
                        <Route exact path='/' render={() => (
                            <Status currentPolicy={props.currentPolicy} isDarkTheme={darkTheme} dt={dt}/>
                        )} />
                        <Route exact path='/details' render={() => (
                            <Details isDarkTheme={darkTheme} dt={dt} errorSetter={props.errorSetter}/>
                        )} />
                    </Box>
                        :
                        <div>
                            <CircularProgress  />
                        {
                            dt.notStarted &&
                            <div>
                                <Typography>Competition have not started yet!</Typography>
                                <Typography>This window will automatically reload once the first round is scored.</Typography>
                            </div>
                        }
                        </div>
                    }
                    {handleFullScreen.active ?
                        <Button
                            style={fullScreenButtonStyle}
                            startIcon={<FullscreenExitIcon />}
                            onClick={handleFullScreen.exit}
                            onMouseOver={handleFadeOver}
                            onMouseOut={handleFadeOut}
                            className={visible ? classes.visible : classes.hidden}

                        >Exit Full Screen</Button>
                        :
                        <Button
                            style={fullScreenButtonStyle}
                            startIcon={<FullscreenIcon />}
                            onClick={handleFullScreen.enter}
                            onMouseOver={handleFadeOver}
                            onMouseOut={handleFadeOut}
                            className={visible ? classes.visible : classes.hidden}
                        >Enter Full Screen</Button>
                    }
            </Box>
        </Paper>
    );
}