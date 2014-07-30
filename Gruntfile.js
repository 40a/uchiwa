'use strict';

module.exports = function (grunt) {
  require('load-grunt-tasks')(grunt);
  require('time-grunt')(grunt);

  grunt.initConfig({
    settings: {
      entryPoint: 'app.js',
      lib: 'lib',
      public: 'public'
    },
    jshint: {
      options: {
        jshintrc: '.jshintrc',
        reporter: require('jshint-stylish')
      },
      all: [
        'Gruntfile.js',
        '<%= settings.entryPoint %>',
        '<%= settings.lib %>/{,*/}*.js',
        '<%= settings.public %>/js/{,*/}*.js',
      ]
    },
    karma: {
      unit: {
        configFile: 'test/karma.conf.js'
      }
    },
    sass: {
      dist: {
        files: [{
          expand: true,
          cwd: '<%= settings.public %>/css',
          src: ['**/*.scss'],
          dest: '<%= settings.public %>/css',
          ext: '.css'
        }]
      }
    },
    watch: {
      scss: {
        files: ['<%= settings.public %>/css/**/*.scss'],
        tasks: ['sass:dist'],
        options: {
          spawn: false
        }
      }
    }
  });

  grunt.registerTask('dev', [
    'watch:scss'
  ]);

  grunt.registerTask('lint', [
    'newer:jshint'
  ]);

  grunt.registerTask('default', [
    'sass',
    'karma:unit'
  ]);
};
