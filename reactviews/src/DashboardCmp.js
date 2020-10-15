import React, { Component } from 'react';
import { MyConfig } from './MyConfig';
import request from 'superagent';
import jwt_decode from 'jwt-decode';

import withWidth from '@material-ui/core/withWidth';
import { AppBar, Title, Responsive } from 'react-admin';

import { DashTables } from './DashTables';
import { DashText } from './DashText';

const styles = {
    docs: { 
        icon: { float: 'right', width: 48, height: 48, padding: 16, color: '#2196F3' },
        card: { borderLeft: 'solid 4px #2196F3', flex: 1, marginLeft: '1em' }
    },
    agents: { 
        icon: { float: 'right', width: 48, height: 48, padding: 16, color: '#4caf50' },
        card: { borderLeft: 'solid 4px #4caf50', flex: 1, marginLeft: '1em' }
    },
    flex: { display: 'flex', marginBottom: '1em' },
    leftCol: { flex: 5, marginRight: '1em' },
    rightCol: { flex: 4, marginLeft: '1em' },
    singleCol: { marginTop: '2em' },
};

class DashboardCmp extends Component {
    state = {};

    handleGetBoard() {
        const token = localStorage.getItem('token');
        request
            .get(MyConfig.API_URL + '/admin/api/v1/board')
            //.set('X-MyToken', MyConfig.API_KEY)
            .set('Authorization', `Bearer ${token}`)
            .end(function(error, response){
                if ( (error && error.status === 401 ) || response === undefined) {
                    localStorage.removeItem('token');
                    window.location.href = MyConfig.BASE_PATH + '#/login';
                    return;
                }
                var r = response.body;
                this.setState({
                    record: r,
                    nbAgents: r.agents,
                    nbSurveys: r.surveys,
                    nbAlerts: r.alerts
                });

            }.bind(this));

    }

    constructor(props) {
        super(props);
        this.handleGetBoard = this.handleGetBoard.bind(this);
        this.tick = this.tick.bind(this);
        //this.clearInterval =  this.clearInterval.bind(this);
    }


    componentDidMount() {
        this.handleGetBoard();
        let timer = setInterval(this.tick, 60000);
        const token = localStorage.getItem('token');
        if (token !== null) {
            var decoded = jwt_decode(localStorage.getItem('token'));
            this.setState({username: decoded.id});
        }
        this.setState({timer});
    }

    componentWillUnmount() {
        clearInterval(this.state.timer);
    } 

    tick() {
        this.handleGetBoard();
    }

    render(props) {
        const {
            nbAgents,
            nbSurveys,
            nbAlerts,
            username,
            record,
        } = this.state;
        return (
            <Responsive
                small = {
                    <div>
                        <AppBar title="DIOC" />
                        <div style={styles.flex}>
                            <DashTables nbagents={nbAgents} nbsurveys={nbSurveys} nbalerts={nbAlerts} subtitle={username}  />
                        </div>
                        <div style={styles.flex}>
                            <DashText record={record} 
                                style={styles.docs} 
                                name="Searchs" 
                                subtitle="Tips and tricks"
                                txt="docsearchs"  />
                        </div>
                        <div style={styles.flex}>
                            <DashText record={record} 
                                style={styles.docs} 
                                name="Doc" 
                                subtitle="Main doc"
                                txt="docs"  />
                        </div>
                    </div>
                }
                medium={
                    <div>
                        <Title title="Distributed IOC manager" />
                        <div style={styles.flex}>
                            <div style={styles.leftCol}>
                                <div style={styles.flex}>
                                    <DashTables nbagents={nbAgents} nbsurveys={nbSurveys} nbalerts={nbAlerts} subtitle={null}  />
                                </div>
                                <div style={styles.flex}>
                                    <DashText record={record} 
                                        style={styles.agents} 
                                        name="Agents" 
                                        subtitle="Downloadable agents"
                                        txt="docagents"  />
                                </div>
                                <div style={styles.flex}>
                                    <DashText record={record} 
                                        style={styles.docs} 
                                        name="Searchs" 
                                        subtitle="Tips and tricks"
                                        txt="docsearchs"  />
                                </div>
                            </div>
                            <div style={styles.rightCol}>
                                <div style={styles.flex}>
                                    <DashText record={record} 
                                        style={styles.docs} 
                                        name="Doc" 
                                        subtitle="Main doc"
                                        txt="docs"  />
                                </div>
                            </div>
                        </div>
                    </div>
                }
            />
        );
    }

}
export default withWidth()(DashboardCmp);
