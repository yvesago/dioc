import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Button } from 'react-admin';
import { showNotification as showNotificationAction } from 'react-admin';
import { push as pushAction } from 'connected-react-router';
import request from 'superagent';
import { MyConfig } from './MyConfig';


class ActionExtractButton extends Component {
    handleClick = () => {
        const token = localStorage.getItem('token');
        const { push, showNotification } = this.props;
        request
            .put(MyConfig.API_URL + '/admin/api/v1/actionextract') 
            .set('Authorization', `Bearer ${token}`)
            .then( (response)  => { 
                showNotification('Extracts done'); 
                /*console.log(response); */
                push('/ips'); 
            })
            .catch((e) => {
                // console.error(e);
                showNotification('Error in action extract ', 'warning');
            });

    }

    render() {
        return <Button label="Extract" onClick={this.handleClick} />;
    }
}

ActionExtractButton.propTypes = {
    push: PropTypes.func,
    showNotification: PropTypes.func,
};


export default connect(null, {
    showNotification: showNotificationAction,
    push: pushAction,
})(ActionExtractButton);

