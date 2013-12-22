#include "slider.h"
#include <cmath>
#include <cstdlib>
#include <QProcess>
#include <QDebug>

Slider::Slider() {
    if(!brightness.open(QIODevice::ReadWrite | QIODevice::Text)) {
        qDebug() << "could not open "<< brightness.fileName();
        std::exit(EXIT_FAILURE);
    }

    QFile maxBrightness{ maxBrightnessPath };
    if(!maxBrightness.open(QIODevice::ReadOnly | QIODevice::Text)) {
        qDebug() << "could not open "<< maxBrightness.fileName();
        std::exit(EXIT_FAILURE);
    }

    char buffer[10] = {};

    maxBrightness.readLine(buffer, sizeof buffer);
    maxBrightnessValue = std::atoi(buffer);
    maxBrightness.close();

    brightness.readLine(buffer, sizeof buffer);
    double brightnessValue = std::atoi(buffer);

    engine.load(QUrl{"qrc:/main.qml"});

    auto backlight = engine.rootObjects().first();
    connect(backlight, SIGNAL(onSlide(qreal)), this, SLOT(onSlide(qreal)));
    backlight->setProperty("slideValue", std::pow(brightnessValue / maxBrightnessValue, 0.5));
}

void Slider::onSlide(qreal value) {
    brightness.resize(0);
    auto string_value = std::to_string(static_cast<int>(std::pow(value, 2) * maxBrightnessValue));
    brightness.write(string_value.c_str(), string_value.size());
    brightness.flush();
}
