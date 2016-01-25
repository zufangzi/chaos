var myApp = angular.module('myApp', ['ngRoute', 'datatables', 'ngResource', 'datatables.bootstrap', 'ui.router']);

myApp.config(function($stateProvider, $urlRouterProvider) {
    $stateProvider
        .state('index', {
            url: "",
            templateUrl: "/static/html/main.html"
        })
        .state('global', {
            url: "/globalinfo",
            templateUrl: "/static/html/global/globalinfo.html"
        })
        .state('apps', {
            url: "/apps",
            templateUrl: "/static/html/apps/apps.html"
        })
        .state('apps.clusterInfo', { // 注意，这边是递进的，url会自动加上前缀
            url: "/clusterinfo",
            templateUrl: "/static/html/apps/apps_cluster_info.html",
            controller: "clusterInfoCtrl"
        });
});

myApp.controller('clusterInfoCtrl', function($scope, $resource, DTOptionsBuilder, DTColumnBuilder) {
    console.log("cluster page fresh...")
        // var vm = this;
    $scope.data = $resource('/api/info').query()
    $scope.dtOptions = DTOptionsBuilder.newOptions()
        .withBootstrap()
        .withBootstrapOptions({
            TableTools: {
                classes: {
                    container: 'btn-group',
                    buttons: {
                        normal: 'btn btn-danger'
                    }
                }
            },
            ColVis: {
                classes: {
                    masterButton: 'btn btn-primary'
                }
            },
            pagination: {
                classes: {
                    ul: 'pagination pagination-sm'
                }
            }
        })
        .withLanguage({
            "sEmptyTable": "亲～表格没有数据",
            "sInfo": "本页 _START_ 到 _END_ 共 _TOTAL_ 条数据",
            "sInfoEmpty": " 0 到 0 共 0 条数据",
            "sInfoFiltered": "(从 _MAX_ total 中筛选)",
            "sInfoPostFix": "",
            "sInfoThousands": ",",
            "sLengthMenu": "每页显示 _MENU_ 条记录",
            "sLoadingRecords": "加载中...",
            "sProcessing": "处理中...",
            "sSearch": "搜索:",
            "sZeroRecords": "没有找到符合条件的数据",
            "oPaginate": {
                "sFirst": "首页",
                "sLast": "尾页",
                "sNext": "下页",
                "sPrevious": "上页"
            },
            "oAria": {
                "sSortAscending": ": 对该列进行升序",
                "sSortDescending": ": 对该列进行降序"
            }
        });
});