import argparse
# import logging
import os
import re
import yaml


class Chart:
    def __init__(self, api_version, name, version):
        self.api_version = api_version
        self.name = name
        self.version = version


class AppConfiguration:
    def __init__(self, config_version, metadata, entrance, permission, spec):
        self.config_version = config_version
        self.metadata = metadata
        self.entrance = entrance
        self.permission = permission
        self.spec = spec


class AppMetaData:
    def __init__(self, name, icon, description, appid, title, version, categories):
        self.name = name
        self.icon = icon
        self.description = description
        self.appid = appid
        self.title = title
        self.version = version
        self.categories = categories


def validate_chart_folder(folder):
    # Check if the folder name is valid
    if not is_valid_folder_name(folder):
        raise ValueError(f"Invalid folder name: '{folder}' must match '^[a-z0-9]{1,30}$'")

    if not os.path.exists(folder):
        raise FileNotFoundError(f"Folder does not exist: '{folder}'")

    # Check if Chart.yaml file exists
    chart_file = os.path.join(folder, "Chart.yaml")
    if not os.path.isfile(chart_file):
        raise FileNotFoundError(f"Missing Chart.yaml in folder: '{folder}'")

    # Read and parse Chart.yaml file
    with open(chart_file, "r") as file:
        chart_content = file.read()
    chart_data = yaml.safe_load(chart_content)
    chart = Chart(
        api_version=chart_data.get("apiVersion"),
        name=chart_data.get("name"),
        version=chart_data.get("version")
    )

    # Check if Chart.yaml fields are valid
    if not is_valid_chart_fields(chart):
        raise ValueError(f"Invalid fields in Chart.yaml in folder: '{folder}'")

    # Check if values.yaml file exists
    values_file = os.path.join(folder, "values.yaml")
    if not os.path.isfile(values_file):
        raise FileNotFoundError(f"Missing values.yaml in folder: '{folder}'")

    # Check if templates folder exists
    templates_dir = os.path.join(folder, "templates")
    if not os.path.isdir(templates_dir):
        raise FileNotFoundError(f"Missing templates folder in folder: '{folder}'")

    # Check if app.cfg file exists
    app_cfg_file = os.path.join(folder, "app.cfg")
    if not os.path.isfile(app_cfg_file):
        raise FileNotFoundError(f"Missing app.cfg in folder: '{folder}'")

    # Read and parse app.cfg file
    with open(app_cfg_file, "r") as file:
        app_cfg_content = file.read()
    app_cfg_data = yaml.safe_load(app_cfg_content)
    app_conf = AppConfiguration(
        config_version=app_cfg_data.get("app.cfg.version"),
        metadata=AppMetaData(**app_cfg_data.get("metadata")),
        entrance=app_cfg_data.get("entrance"),
        permission=app_cfg_data.get("permission"),
        spec=app_cfg_data.get("spec")
    )

    # Check if metadata fields in app.cfg are valid
    error_message = is_valid_metadata_fields(app_conf.metadata, chart, folder)
    if error_message:
        raise ValueError(f"Invalid metadata fields in app.cfg in folder: '{folder}'. {error_message}")

    # if not is_valid_metadata_fields(app_conf.metadata, chart, folder):
    #     raise ValueError(f"Invalid metadata fields in app.cfg in folder: '{folder}'")

    # Validation passed
    return True


def is_valid_folder_name(name):
    match = re.match("^[a-z0-9]{1,30}$", name)
    return bool(match)


def is_valid_chart_fields(chart):
    if not chart.api_version:
        return False
    if not chart.name:
        return False
    if not chart.version:
        return False
    return True

def is_valid_metadata_fields(metadata, chart, folder):
    if chart.name != folder:
        return f"name {chart.name} invalid in Chart.yaml in chart '{folder}', must be same"

    if metadata.name != folder:
        return f"metadata.name {metadata.name} invalid in app.cfg in chart '{folder}', must be same"

    if metadata.version != chart.version:
        return f"version in app.cfg {metadata.version}, version in Chart.yaml {chart.version} in chart '{folder}', must be same"

    return None


# def is_valid_metadata_fields(metadata, chart, folder):
#     if chart.name != folder:
#         return False
#     if metadata.name != folder:
#         return False
#     if metadata.version != chart.version:
#         return False
#     return True


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-chart-dirs", help="comma-separated list of chart directories")
    args = parser.parse_args()

    if not args.chart_dirs:
        parser.print_usage()
        return

    dirs = args.chart_dirs.split(",")
    for dir in dirs:
        try:
            validate_chart_folder(dir)
        except Exception as e:
            # logging.exception(e)
            print(f"Validation failed for folder '{dir}': {str(e)}")
            return

    print("Folder validation successful!")


if __name__ == "__main__":
    main()