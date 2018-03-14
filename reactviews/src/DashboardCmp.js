import React, { Component } from 'react';
import { MyConfig } from './MyConfig';
import request from 'superagent';

import withWidth from 'material-ui/utils/withWidth';
import { AppBarMobile } from 'admin-on-rest';

import DashTables from './DashTables';
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
        request
            .get(MyConfig.API_URL + '/admin/api/v1/board')
            .set('X-MyToken', MyConfig.API_KEY)
            .end(function(error, response){
                if ( (error && error.status === 401 ) || response === undefined) {
                    window.location.href = '/#/';
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
            record,
        } = this.state;
        const { width } = this.props;
        return (
            <div>
                {width === 1 && <AppBarMobile title="DashBoard" />}
                <div style={styles.flex}>
                    <div style={styles.leftCol}>
                        <div style={styles.flex}>
                            <DashTables nbagents={nbAgents} nbsurveys={nbSurveys} nbalerts={nbAlerts}  />
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
        );
    }

}
export default withWidth()(DashboardCmp);
