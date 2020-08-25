import React from "react";
import AuthService from "../../services/auth/auth";
import SingleTeamDetails from "./SingleTeamDetails";
import TableContainer from "@material-ui/core/TableContainer";
import Paper from "@material-ui/core/Paper";
import TableHead from "@material-ui/core/TableHead";
import Table from "@material-ui/core/Table";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import TableBody from "@material-ui/core/TableBody";
import IconButton from "@material-ui/core/IconButton";
import { makeStyles } from '@material-ui/core/styles';
import Box from '@material-ui/core/Box';
import Collapse from '@material-ui/core/Collapse';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';

export default function Details(props) {
    const dt=props.dt

    function BlackTeamPanel() {
        let data = []
        Object.keys(dt["Teams"]).forEach(team_id =>{
            data.push({
                team_id: team_id,
                team_name: dt["Teams"][team_id]["Name"],
            })
        })

        data.sort((a, b) => (a.team_name > b.team_name) ? 1 : -1)

        return (

            <TableContainer component={Paper}>
                <Table aria-label="collapsible table">
                    <TableHead>
                        <TableRow>
                            <TableCell />
                            <TableCell align="right">Team Name</TableCell>
                            <TableCell align="right">Team ID</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {data.map((row) => (
                            <CustomRow key={row.team_id} {...props} row={row} />
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>

        )
    }
    return (
        <div>
            <link
                rel="stylesheet"
                href="https://fonts.googleapis.com/icon?family=Material+Icons"
            />


            {
                AuthService.getCurrentRole() === "blue" ? <SingleTeamDetails {...props} teamID={AuthService.getCurrentTeamID()}/> :
                    BlackTeamPanel()
            }
        </div>
    );
}

const useRowStyles = makeStyles({
    root: {
        '& > *': {
            borderBottom: 'unset',
        },
    },
});



function CustomRow(props) {
    const { row } = props;
    const [open, setOpen] = React.useState(false);
    const classes = useRowStyles();
    return (
        <React.Fragment>
            <TableRow className={classes.root}>
                <TableCell>
                    <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell align="right">{row.team_name}</TableCell>
                <TableCell align="right">{row.team_id}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
                    <Collapse in={open} timeout="auto">
                        <Box margin={1}>
                            <SingleTeamDetails {...props} teamID={row.team_id} />
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}