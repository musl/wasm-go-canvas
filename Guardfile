clearing :on
interactor :on
notification :off

guard :shell do
    ignore(/build\/.*/)
    watch(/.*\.(go|js|html|css)/) { `make clean build` }
end
