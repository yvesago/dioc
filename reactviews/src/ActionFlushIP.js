import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button } from 'react-admin';
import { showNotification as showNotificationAction } from 'react-admin';
import { push as pushAction } from 'react-router-redux';
import request from 'superagent';
import { MyConfig } from './MyConfig';


class ActionFlushIPButton extends Component {
    handleClick = () => {
        const token = localStorage.getItem('token');
        const { push, showNotification } = this.props;
        request
            .put(MyConfig.API_URL + '/admin/api/v1/actionfluship') 
            .set('Authorization', `Bearer ${token}`)
            .then( (response)  => { 
                showNotification('Flush IPs done'); 
                /*console.log(response); */
                push('/extracts'); 
            })
            .catch((e) => {
                // console.error(e);
                showNotification('Error in action Flush IPs', 'warning');
            });

    }

    render() {
        return <Button color="primary" onClick={this.handleClick}>Flush IPs</Button>;
    }
}

ActionFlushIPButton.propTypes = {
    push: PropTypes.func,
    showNotification: PropTypes.func,
};


export default connect(null, {
    showNotification: showNotificationAction,
    push: pushAction,
})(ActionFlushIPButton);

