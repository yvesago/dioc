
var myApp = angular.module('myApp', ['restangular','ng-admin']);
myApp.config(['NgAdminConfigurationProvider','RestangularProvider', function (nga,rp) {

    rp.setDefaultHeaders({"X-MyToken":'basic'}); // XXX fix token

    // create an admin application
    var admin = nga.application('Distributed IOC Monitor')
      .baseApiUrl('http://localhost:8080/admin/api/v1/'); // XXX fix main API endpoint

    var roles = [ // XXX fix your categorys
        { label: 'Squid', value: 'squid' },
        { label: 'Radius', value: 'radius' },
        { label: 'Web', value: 'web' },
        { label: 'DNS', value: 'dns' },
        { label: 'Honeypot', value: 'honeypot' },
        { label: 'NetFlow', value: 'netflow' },
        { label: 'Auth', value: 'auth' },
        { label: 'SMTP', value: 'smtp' },
        { label: 'Mail', value: 'mail' },
        { label: 'Test', value: 'test' },
        { label: 'None', value: '' }
    ];

    var levels = [
        { label: 'Warn', value: 'warn' },
        { label: 'Critic', value: 'critic' },
        { label: 'None', value: '' }
    ];

    // create a survey entity
    var survey = nga.entity('surveys');
    survey.listView().fields([
        nga.field('search').isDetailLink(true),
        nga.field('role'),
        nga.field('level').cssClasses(['red']).template('<span ng-class="{ \'orange\': value === \'warn\' }"><ma-string-column field="::field" value="::entry.values[field.name()]"></ma-string-column></span>'),
        nga.field('crcs'),
        nga.field('comment','wysiwyg'),
        nga.field('created','datetime'),
        nga.field('updated','datetime')
    ]).filters([
            nga.field('search')
            .label('Search')
            .pinned(true),
            nga.field('comment')
            .label('Comment')
            ]);
    survey.creationView().fields([
        nga.field('search').validation({ required: true }),
        nga.field('level', 'choice').choices(levels),
        nga.field('role', 'choice').choices(roles),
        nga.field('comment','wysiwyg')
    ]);
    survey.editionView().fields([
        nga.field('search').validation({ required: true }),
        nga.field('crcs').editable(false),
        nga.field('created','datetime').editable(false),
        nga.field('updated','datetime').editable(false),
        nga.field('level', 'choice').choices(levels),
        nga.field('role', 'choice').choices(roles),
        nga.field('comment','wysiwyg')
    ]);
    admin.addEntity(survey);

    var agent = nga.entity('agents').identifier(nga.field('crca'));
    agent.listView().fields([
    //   nga.field('crca').isDetailLink(true),
       nga.field('ip').isDetailLink(true),
       nga.field('filesurvey'),
       nga.field('status').cssClasses(['red','green']).template('<span ng-class="{ \'red\': value === \'OffLine\' }"><ma-string-column field="::field" value="::entry.values[field.name()]"></ma-string-column></span>'),
       nga.field('role'),
       nga.field('lines','text'),
       nga.field('comment','text'),
       nga.field('cmd'),
       nga.field('created','datetime'),
       nga.field('updated','datetime')
    ]).filters([
            nga.field('ip')
            .label('IP')
            .pinned(true),
            nga.field('comment')
            .label('Comment'),
            nga.field('lines')
            .label('Lines')
            ]);
    agent.editionView().fields([
       nga.field('crca').editable(false),
       nga.field('created','datetime').editable(false),
       nga.field('updated','datetime').editable(false),
       nga.field('status').cssClasses(['red','green']).template('<span ng-class="{ \'red\': value === \'OffLine\' }"><ma-string-column field="::field" value="::entry.values[field.name()]"></ma-string-column></span>'),
       nga.field('ip'),
       nga.field('filesurvey'),
       nga.field('lines','text'),
       nga.field('comment','text'),
       nga.field('role', 'choice').choices(roles),
       nga.field('cmd', 'choice').choices([
                       { label: 'Send Lines', value: 'SendLines' },
                       { label: 'Search Full File', value: 'FullSearch' },
                       { label: 'STOP', value: 'STOP' },
                       { label: 'None', value: '' }
                   ])
    ]);
    //agent.editionView().fields(agent.creationView().fields());
    admin.addEntity(agent);

    var alerte = nga.entity('alertes');
    alerte.listView().fields([
        nga.field('search').isDetailLink(true),
        nga.field('role'),
        nga.field('ip'),
        nga.field('filesurvey'),
        nga.field('level').cssClasses(['red']).template('<span ng-class="{ \'orange\': value === \'warn\' }"><ma-string-column field="::field" value="::entry.values[field.name()]"></ma-string-column></span>'),
        nga.field('line'),
        nga.field('comment','text'),
        nga.field('created','datetime'),
        nga.field('updated','datetime')
    ]).filters([
            nga.field('line')
            .label('Line')
            .pinned(true),
            nga.field('comment')
            .label('Comment'),
            nga.field('line')
            .label('Line')
            ]);
    alerte.editionView().fields([
        nga.field('crca').editable(false),
        nga.field('crcs').editable(false),
        nga.field('ip').editable(false),
        nga.field('filesurvey').editable(false),
        nga.field('role').editable(false),
        nga.field('level', 'choice').choices(levels),
        nga.field('line'),
        //nga.field('comment','text')
        nga.field('comment','wysiwyg'),
    ]);
    //alerte.editionView().fields(alerte.creationView().fields());
    admin.addEntity(alerte);
	admin.menu(nga.menu()
            .addChild(nga.menu(survey))
            .addChild(nga.menu(agent))
            .addChild(nga.menu(alerte))
	    // add custom menu
            .addChild(nga.menu().title('Miscellaneous').icon('<span class="glyphicon glyphicon-wrench"></span>')
              .addChild(nga.menu().title('Download Agent').link('/dwnloadagent').icon('<span class="glyphicon glyphicon-file"></span>'))
	      .addChild(nga.menu().template(`
		<a href="https://github.com/yvesago/dioc" target="_blank">
		    <span class="glyphicon glyphicon-download"></span>
		    Source code
		</a>`
	      ))
            )
	);

    // attach the admin application to the DOM and execute it
    nga.configure(admin);
}]);

var dwnloadAgentTemplate =
    '<div class="row"><div class="col-lg-12">' +
        '<ma-view-actions><ma-back-button></ma-back-button></ma-view-actions>' +
        '<div class="page-header">' +
            '<h1>Download agent</h1>' +
        '</div>' +
    '</div></div>' +
    '<div class="row">' +
        '<div class="col-lg-10"><tt>[MD5sum] <a href="/ioc/agent" target="_blank">agent</a></tt></div>' +
        '<div class="col-lg-10"><tt>[MD5sum] <a href="/ioc/start-stop-agent" target="_blank">start-stop-agent</a></tt></div>' +
    '</div>';

myApp.config(function ($stateProvider) {
    $stateProvider.state('dwnloadagent', {
	parent: 'ng-admin',
	url: '/dwnloadagent',
	template: dwnloadAgentTemplate
    });
});
