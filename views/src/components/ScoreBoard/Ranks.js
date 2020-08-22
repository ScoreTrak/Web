import { ResponsiveBar } from '@nivo/bar'

import React from "react";

const darkTheme = {
    axis: {
        fontSize: "14px",
        tickColor: "#eee",
        ticks: {
            line: {
                stroke: "#555555"
            },
            text: {
                fill: "#ffffff"
            }
        },
        legend: {
            text: {
                fill: "#aaaaaa"
            }
        }
    },
    grid: {
        line: {
            stroke: "#555555"
        }
    },
    tooltip: {
        container: {
            background: '#333',
        },
    },
};


export default function Ranks(props) {
    const dt=props.dt
    let data = []
    let dataKeys = new Set();
    if ("Teams" in dt){
        for (let team in dt["Teams"]) {
            if (dt["Teams"].hasOwnProperty(team)) {
                let serviceAggregator = {}
                for (let host in dt.Teams[team]["Hosts"]){
                    if (dt.Teams[team]["Hosts"].hasOwnProperty(host)) {
                        if (Object.keys(dt.Teams[team]["Hosts"][host]["Services"]).length !== 0){
                            for (let service in dt.Teams[team]["Hosts"][host]["Services"]) {
                                if (dt.Teams[team]["Hosts"][host]["Services"].hasOwnProperty(service)) {
                                    let sr = dt.Teams[team]["Hosts"][host]["Services"][service]
                                    let keyName = ""
                                    if (sr["DisplayName"]){
                                        keyName = sr["DisplayName"]
                                    } else {
                                        if (dt.Teams[team]["Hosts"][host]["HostGroup"]){
                                            keyName = dt.Teams[team]["Hosts"][host]["HostGroup"]["Name"] + "-" + sr["Name"]
                                        } else{
                                            keyName = sr["Name"]
                                        }
                                    }
                                    dataKeys.add(keyName)
                                    if (keyName in serviceAggregator){
                                        serviceAggregator[keyName] += sr["Points"] + sr["PointsBoost"]
                                    } else {
                                        serviceAggregator[keyName] = sr["Points"] + sr["PointsBoost"]
                                    }
                                }
                            }

                        }
                    }
                }

                if (Object.keys(serviceAggregator).length !==0){
                    data.push({
                        ...serviceAggregator,
                        teamName: dt["Teams"][team]["Name"],
                    })
                }
            }
        }
    }

    const serviceSum = (teamObj) => {
        let sum = 0
        Object.keys(teamObj).forEach(function(key) {
            if (key !== "teamName"){
                sum += teamObj[key];
            }
        });
        return sum
    }

    data.sort((a, b) => (serviceSum(a) === serviceSum(b)) ? (a.teamName > b.teamName ? 1 : -1) : (serviceSum(a) > serviceSum(b) ? 1 : -1))

    let theme= {fontSize: "0.875rem"}
    if (props.isDarkTheme){
        Object.assign(theme, darkTheme);
    }
    return (
        <ResponsiveBar
            theme={theme}
            data={data}
            keys={[ ...dataKeys ]}
            indexBy="teamName"
            margin={{ top: 50, right: 60, bottom: 50, left: 60 }}
            padding={0.3}
            colors={{ scheme: props.isDarkTheme ? 'nivo' : 'dark2' }}
            borderColor={{ from: 'color', modifiers: [ [ 'darker', '0' ] ] }}
            axisTop={null}
            axisRight={null}
            axisLeft={{
                tickSize: 5,
                tickPadding: 5,
                tickRotation: 0,
                legend: 'points',
                legendPosition: 'middle',
                legendOffset: -40
            }}
            labelSkipWidth={8}
            labelSkipHeight={12}
            labelTextColor={{ from: 'color', modifiers: [ [ 'darker', 2 ] ] }}
            legends={[]}
            animate={true}
            motionStiffness={70}
            motionDamping={15}
        />
    );
}