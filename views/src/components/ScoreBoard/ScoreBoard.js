import React, {useState, useRef, useEffect} from "react";
import ReportService from "../../services/report/report"
import RoundService from "../../services/round/round"
import CircularProgress from '@material-ui/core/CircularProgress';
import Box from '@material-ui/core/Box';
import {Typography} from "@material-ui/core";
import { Route } from "react-router-dom";
import Ranks from "./Ranks";
import Uptime from "./Uptime";
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import Button from '@material-ui/core/Button';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import {makeStyles} from "@material-ui/core/styles";

const fullScreenButtonStyle = {
    position: 'absolute',
    maxHeight: '3vh',
    bottom: '0px',
    marginRight: '2vh',
    right: '20px',
}

const useStyles = makeStyles((theme) => ({
    "hidden": { opacity: 0.1, transition: "opacity 0.2s linear"},
    "visible": { opacity: 1, transition: "opacity 0.2s linear"},
}));

export default function ScoreBoard(props) {
    const [dt, setData] = useState({Round: 0, loader:true});
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
    console.log(dt)
    useEffect(() => {
        function reload() {
            RoundService.getLastNonElapsingRound().then( round =>{
                if (round.id !== roundRef.current){
                    ReportService.getReport().then( simpleReport => {
                        roundRef.current = simpleReport.Round
                        setData(simpleReport)
                        props.setTitle("Round " + simpleReport.Round)
                    })
                }

            },
                (error) => {
                    const resMessage =
                        (error.response &&
                            error.response.data &&
                            error.response.data.message) ||
                        error.message ||
                        error.toString();
                    props.setError(resMessage)
                    setData({loader: true, Round: 0,})
                })
        }
        reload()
        const roundReload = setInterval(reload, 3000);
        return () => clearInterval(roundReload);
    }, []);

    return (
            <Box m="auto" height="100%" width="100%" >
                { !("loader" in dt) && dt.Round !== 0 ?
                    <Box m="auto" height="100%" width="100%">
                        <Route exact path='/' render={() => (
                            <Ranks isDarkTheme={darkTheme} dt={dt}/>
                        )} />
                        <Route exact path='/uptime' render={() => (
                            <Uptime isDarkTheme={darkTheme} dt={dt}/>
                        )} />
                    </Box>
                        :
                        <div>
                            <CircularProgress  />
                        {
                            dt.Round === 0  && !("loader" in dt) &&
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

    );
}