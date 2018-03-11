import React, { Component } from 'react';
import { MyConfig } from './MyConfig';
import request from 'superagent';

import withWidth from 'material-ui/utils/withWidth';
import { AppBarMobile } from 'admin-on-rest';

import DashTables from './DashTables';

const styles = {
    welcome: { marginBottom: '2em' },
    flex: { display: 'flex', marginBottom: '1em' },
    leftCol: { flex: 1, marginRight: '1em' },
    rightCol: { flex: 1, marginLeft: '1em' },
    singleCol: { marginTop: '2em' },
};

class DashboardCmp extends Component {
    state = {};

    handleGetBoard() {
        request
            .get(MyConfig.API_URL + '/admin/api/v1/board')
            .set('X-MyToken', MyConfig.API_KEY)
            .end(function(error, response){
                if (error && error.status === 401) {
                    window.location.href = '/#/';
                    return;
                }
                var r = response.body;
                this.setState({
                    //record: r,
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
            //record,
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
                    </div>
                </div>
            </div>
        );
    }

}
export default withWidth()(DashboardCmp);
