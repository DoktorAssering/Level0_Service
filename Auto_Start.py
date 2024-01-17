import subprocess
import time
import os

current_dir = os.path.dirname(os.path.abspath(__file__))

commands_to_execute = [
    "start cmd /k nats-streaming-server",
    f"cd {os.path.join(current_dir,'ServiceApp', 'nats')} && start cmd /k go run Nats.go",
    f"cd {os.path.join(current_dir,'ServiceApp', 'server')} && start cmd /k go run Server.go"
]

def execute_commands(commands):
    try:
        for command in commands:
            print(f"Executing command: {command}")
            
            subprocess.Popen(command, shell=True)
            
            print("=" * 50)

            time.sleep(2)
        
    except Exception as e:
        print(f"Error executing commands: {e}")

execute_commands(commands_to_execute)
