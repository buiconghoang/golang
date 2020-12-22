#!/usr/bin/env python
# coding: utf-8

# In[58]:


import re


# In[69]:


def is_able_to_contain_linux_paths(command: str):
    pattern = r"(^/|[\s'\"]/)"
    re.compile(pattern)
    r = re.findall(pattern, command)
#     print(r)
    if len(r) > 0:
        return True
    return False

commands = ["/usr/local/bin", "\"/usr/local/bin\"", "move '/abc/abc' ", "move /abc/abc"]
for command in commands:
    print(command, " --- ", is_able_to_contain_linux_paths(command))


# In[71]:


def is_able_to_contain_window_paths(command: str):
    disk_prefix = r"((^[A-Za-z0-9]+|[\s'\"][A-Za-z0-9]+):(\\+|\/+))"
    keep_prefix = r"((^(swam|home)|[\s'\"](swam|home))(\\+|\/+))"
    patterns = [keep_prefix, disk_prefix]
    re.compile(keep_prefix)
    re.compile(disk_prefix)
    
    for pattern in patterns:
        r = re.findall(pattern, command)
        if len(r) > 0:
#             print(r)
            return True
    
    return False

commands = ["C:\\abcde\\ab.py", "abc C:\\abcde\\ab.py",  "\"D://hello\\abc\"", "move 'swam/abc' ", "move \"home\\hello\\a.py\" abc ", "'D://hello\\abc'"]
for command in commands:
    print(command, " --- ", is_able_to_contain_window_paths(command))


# In[74]:


import os
def export_linux_window_path(datapath: str):
    file_names = os.listdir(datapath)
    window_path_writer = open("window_paths.txt", mode='w')
    linux_path_writer = open("linux_paths.txt", mode='w')
    for file_name in file_names:
        file_path = os.path.join(datapath, file_name)
        with open(file_path, 'r') as f:
            line = f.readline()
            while line:
                line = line.split("\n")[0]
                parts = line.split("\t")
#                 print(parts)
                if len(parts) > 1:
                    
                    command = parts[1]
                    if is_able_to_contain_window_paths(command):
                        window_path_writer.write(command + "\n")
                        window_path_writer.flush()
                    elif is_able_to_contain_linux_paths(command):
                        linux_path_writer.write(command + "\n")
                        linux_path_writer.flush()
                        
                line = f.readline()
    
    window_path_writer.close()
    linux_path_writer.close()
    
    
export_linux_window_path("./data")


# In[75]:


import os
os.listdir("./data")


# In[76]:


os.getcwd()


# In[ ]:




